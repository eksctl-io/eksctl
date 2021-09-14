package credentials

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gofrs/flock"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	// EksctlGlobalEnableCachingEnvName defines an environment property to enable the cache globally.
	EksctlGlobalEnableCachingEnvName = "EKSCTL_ENABLE_CREDENTIAL_CACHE"
	// EksctlCacheFilenameEnvName defines an environment property to configure where the cache file should live.
	EksctlCacheFilenameEnvName = "EKSCTL_CREDENTIAL_CACHE_FILENAME"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
// Clock implements Now to return the current time.
//counterfeiter:generate -o fakes/fake_clock.go . Clock
type Clock interface {
	Now() time.Time
}

// RealClock defines a clock using time.Now()
type RealClock struct{}

// Now returns the current time.
func (r *RealClock) Now() time.Time {
	return time.Now()
}

type cachedCredential struct {
	Credential credentials.Value
	Expiration time.Time
}

// FileCacheProvider is a file based AWS Credentials Provider implementing expiry and retrieve.
type FileCacheProvider struct {
	credentials      *credentials.Credentials // the underlying implementation that has the *real* Provider
	cachedCredential cachedCredential         // the cached credential, if it exists
	profile          string
	clock            Clock
}

type cacheFile struct {
	// a map of profiles to cachedCredentials
	ProfileMap map[string]cachedCredential `yaml:"profiles"`
}

// Put puts the given cachedCredential with a given key into the map. It will overwrite
// if the key already exists.
func (c *cacheFile) Put(key string, credential cachedCredential) {
	c.ProfileMap[key] = credential
}

// Get returns cachedCredential if it exists in the cred store.
func (c *cacheFile) Get(key string) cachedCredential {
	var credential cachedCredential
	if _, ok := c.ProfileMap[key]; ok {
		credential = c.ProfileMap[key]
	}
	return credential
}

// NewFileCacheProvider creates a new filesystem based AWS credential cache. The cache uses Expiry provided by the
// AWS Go SDK for providers. It wraps the configured credential provider into a file based cache provider. If the provider
// does not support caching ( I.e.: it doesn't implement IsExpired ) then this file based caching system is ignored
// and the default credential provider is used. Caches are per profile.
func NewFileCacheProvider(profile string, creds *credentials.Credentials, clock Clock) (FileCacheProvider, error) {
	if creds == nil {
		return FileCacheProvider{}, errors.New("no underlying Credentials object provided")
	}
	filename, err := cacheFilename()
	if err != nil {
		return FileCacheProvider{}, fmt.Errorf("failed to get cache file: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return FileCacheProvider{}, fmt.Errorf("failed to create folder: %w", err)
	}
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		logger.Warning("Cache file %s does not exist.\n", filename)
		return FileCacheProvider{
			profile:          profile,
			credentials:      creds,
			cachedCredential: cachedCredential{},
			clock:            clock,
		}, nil
	}

	if info.Mode()&0077 != 0 {
		// cache file has secret credentials and should only be accessible to the user, refuse to use it.
		return FileCacheProvider{}, fmt.Errorf("cache file %s is not private", filename)
	}

	cache, err := readCacheFile(filename)
	if err != nil {
		return FileCacheProvider{}, err
	}

	return FileCacheProvider{
		credentials:      creds,
		cachedCredential: cache.Get(profile),
		profile:          profile,
		clock:            clock,
	}, nil
}

// readCacheFile reads the contents of the credential cache and returns the
// parsed yaml as a cachedCredential object.
func readCacheFile(filename string) (cacheFile, error) {
	lock := flock.New(filename)
	defer func() {
		if err := lock.Unlock(); err != nil {
			logger.Warning("Unable to unlock file %s: %v\n", filename, err)
		}
	}()
	// wait up to a second for the file to lock
	cache := cacheFile{
		ProfileMap: make(map[string]cachedCredential),
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	ok, err := lock.TryRLockContext(ctx, 250*time.Millisecond) // try to lock every 1/4 second
	if !ok {
		// unable to lock the cache, something is wrong, refuse to use it.
		return cache, fmt.Errorf("unable to read lock file %s: %v", filename, err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return cache, fmt.Errorf("failed to read cache file: %w", err)
	}
	if err := yaml.Unmarshal(data, &cache); err != nil {
		return cache, fmt.Errorf("unable to parse file %s: %w", filename, err)
	}

	return cache, nil
}

// writeCache writes the contents of the credential cache using the
// yaml marshaled form of the passed cachedCredential object.
func writeCache(filename string, cache cacheFile) error {
	lock := flock.New(filename)
	defer func() {
		if err := lock.Unlock(); err != nil {
			logger.Warning("Unable to unlock file %s: %v\n", filename, err)
		}
	}()
	// wait up to a second for the file to lock
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	ok, err := lock.TryRLockContext(ctx, 250*time.Millisecond) // try to lock every 1/4 second
	if !ok {
		// unable to lock the cache, something is wrong, refuse to use it.
		return fmt.Errorf("unable to read lock file %s: %v", filename, err)
	}
	data, err := yaml.Marshal(cache)
	if err == nil {
		// write privately owned by the user
		err = os.WriteFile(filename, data, 0600)
	}
	return err
}

// Retrieve implements the Provider interface, returning the cached credential if is not expired,
// otherwise fetching the credential from the underlying Provider and caching the results on disk
// with an expiration time.
func (f *FileCacheProvider) Retrieve() (credentials.Value, error) {
	if !f.cachedCredential.Expiration.Before(f.clock.Now()) {
		// use the cached credential
		return f.cachedCredential.Credential, nil
	}
	logger.Info("No cached credential available.  Refreshing...")
	// fetch the credentials from the underlying Provider
	credential, err := f.credentials.Get()
	if err != nil {
		return credential, err
	}
	expiration, err := f.credentials.ExpiresAt()
	if err != nil {
		// credential doesn't support expiration time, so can't cache, but still return the credential
		logger.Warning("Unable to cache credential: %v\n", err)
		return credential, nil
	}
	// underlying provider supports Expirer interface, so we can cache
	filename, err := cacheFilename()
	if err != nil {
		return credential, err
	}
	f.cachedCredential = cachedCredential{
		Credential: credential,
		Expiration: expiration,
	}
	// overwrite whatever was there before. we don't care about multiple creds for various clusters.
	// if user switches to another role and another profile they have to re-authenticate.
	cache, _ := readCacheFile(filename)
	cache.Put(f.profile, f.cachedCredential)
	if err := writeCache(filename, cache); err != nil {
		logger.Warning("Unable to update credential cache %s: %v\n", filename, err)
		return credential, err
	}
	logger.Info("Updated cached credential\n")
	return credential, nil
}

// IsExpired implements the Provider interface, deferring to the cached credential first,
// but fall back to the underlying Provider if it is expired.
func (f *FileCacheProvider) IsExpired() bool {
	return f.cachedCredential.Expiration.Before(f.clock.Now()) && f.credentials.IsExpired()
}

// ExpiresAt implements the Expirer interface, and gives access to the expiration time of the credential
func (f *FileCacheProvider) ExpiresAt() time.Time {
	return f.cachedCredential.Expiration
}

func cacheFilename() (string, error) {
	if filename := os.Getenv(EksctlCacheFilenameEnvName); filename != "" {
		return filename, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filename := filepath.Join(home, ".eksctl", "cache", "credentials.yaml")
	return filename, nil
}
