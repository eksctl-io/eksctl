package credentials

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/kris-nova/logger"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const (
	// EksctlGlobalEnableCachingEnvName defines an environment property to enable the cache globally.
	EksctlGlobalEnableCachingEnvName = "EKSCTL_ENABLE_CREDENTIAL_CACHE"
	// EksctlCacheFilenameEnvName defines an environment property to configure where the cache file should live.
	EksctlCacheFilenameEnvName = "EKSCTL_CREDENTIAL_CACHE_FILENAME"
)

// Clock implements Now to return the current time.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
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

// Flock provides an interface to handle file locking.
// It defines an interface for the Flock type from github.com/gofrs/flock.
// Refer to https://pkg.go.dev/github.com/gofrs/flock?utm_source=godoc#Flock for documentation.
//
//counterfeiter:generate -o fakes/fake_flock.go . Flock
type Flock interface {
	// TryRLockContext repeatedly tries to take a shared lock until one of the
	// conditions is met: TryRLock succeeds, TryRLock fails with error, or Context
	// Done channel is closed.
	TryRLockContext(ctx context.Context, retryDelay time.Duration) (bool, error)

	// TryLockContext repeatedly tries to take an exclusive lock until one of the
	// conditions is met: TryLock succeeds, TryLock fails with error, or Context
	// Done channel is closed.
	TryLockContext(ctx context.Context, retryDelay time.Duration) (bool, error)

	// Unlock is unlocks the file.
	Unlock() error
}

type FlockFunc func(path string) Flock

type cachedCredential struct {
	Credential credentials.Value
	Expiration time.Time
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

func initializeCache(fs afero.Fs, cacheFilePath string) error {
	if err := fs.MkdirAll(filepath.Dir(cacheFilePath), 0700); err != nil {
		return fmt.Errorf("failed to create folder: %w", err)
	}
	info, err := fs.Stat(cacheFilePath)
	if os.IsNotExist(err) {
		logger.Warning("cache file %s does not exist.\n", cacheFilePath)
		return nil
	}

	if info.Mode()&0077 != 0 {
		// cache file has secret credentials and should only be accessible to the user, refuse to use it.
		return fmt.Errorf("cache file %s is not private", cacheFilePath)
	}

	_, err = parseCacheFile(fs, cacheFilePath)
	return err
}

// readCacheFile reads the contents of the credential cache and returns the
// parsed yaml as a cachedCredential object.
func readCacheFile(fs afero.Fs, filename string, newFlock FlockFunc) (cacheFile, error) {
	cache := cacheFile{
		ProfileMap: make(map[string]cachedCredential),
	}
	if _, err := fs.Stat(filename); os.IsNotExist(err) {
		logger.Warning("cache file %s does not exist.\n", filename)
		return cache, nil
	}
	lock := newFlock(filename)
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
		return cache, fmt.Errorf("unable to read lock file %s: %v", filename, err)
	}
	return parseCacheFile(fs, filename)
}

// writeCache writes the contents of the credential cache using the
// yaml marshaled form of the passed cachedCredential object.
func writeCache(fs afero.Fs, filename string, newFlock FlockFunc, cache cacheFile) error {
	lock := newFlock(filename)
	defer func() {
		if err := lock.Unlock(); err != nil {
			logger.Warning("Unable to unlock file %s: %v\n", filename, err)
		}
	}()
	// wait up to a second for the file to lock
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()
	ok, err := lock.TryLockContext(ctx, 250*time.Millisecond) // try to lock every 1/4 second
	if !ok {
		// unable to lock the cache, something is wrong, refuse to use it.
		return fmt.Errorf("unable to read lock file %s: %v", filename, err)
	}
	data, err := yaml.Marshal(cache)
	if err == nil {
		// write privately owned by the user
		err = afero.WriteFile(fs, filename, data, 0600)
	}
	return err
}

func parseCacheFile(fs afero.Fs, filename string) (cacheFile, error) {
	cache := cacheFile{
		ProfileMap: make(map[string]cachedCredential),
	}
	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return cache, fmt.Errorf("failed to read cache file: %w", err)
	}
	if err := yaml.Unmarshal(data, &cache); err != nil {
		return cache, fmt.Errorf("unable to parse file %s: %w", filename, err)
	}
	return cache, nil
}

// GetCacheFilePath gets the filename to use for caching credentials.
func GetCacheFilePath() (string, error) {
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
