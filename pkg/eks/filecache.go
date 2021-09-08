package eks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type cachedCredential struct {
	Credential credentials.Value
	Expiration time.Time
	// If set will be used by IsExpired to determine the current time.
	// Defaults to time.Now if CurrentTime is not set.  Available for testing
	// to be able to mock out the current time.
	currentTime func() time.Time
}

// IsExpired determines if the cached credential has expired
func (c *cachedCredential) IsExpired() bool {
	curTime := c.currentTime
	if curTime == nil {
		curTime = time.Now
	}
	return c.Expiration.Before(curTime())
}

type FileCacheProvider struct {
	credentials      *credentials.Credentials // the underlying implementation that has the *real* Provider
	cachedCredential cachedCredential         // the cached credential, if it exists
	profile          string
}

type cacheFile struct {
	// a map of profiles to cachedCredentials
	ProfileMap map[string]cachedCredential `yaml:"profiles"`
}

func (c *cacheFile) Put(key string, credential cachedCredential) {
	c.ProfileMap[key] = credential
}

func (c *cacheFile) Get(key string) (credential cachedCredential) {
	if _, ok := c.ProfileMap[key]; ok {
		credential = c.ProfileMap[key]
	}
	return
}

func NewFileCacheProvider(profile string, creds *credentials.Credentials) (FileCacheProvider, error) {
	if creds == nil {
		return FileCacheProvider{}, errors.New("no underlying Credentials object provided")
	}
	filename, err := cacheFilename()
	if err != nil {
		return FileCacheProvider{}, fmt.Errorf("failed to get cache file: %w", err)
	}
	_ = os.MkdirAll(filepath.Dir(filename), 0700)
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Cache file %s does not exist.\n", filename)
		return FileCacheProvider{
			credentials:      creds,
			cachedCredential: cachedCredential{},
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
	}, nil
}

// readCacheFile reads the contents of the credential cache and returns the
// parsed yaml as a cachedCredential object.
func readCacheFile(filename string) (cacheFile, error) {
	cache := cacheFile{
		ProfileMap: make(map[string]cachedCredential),
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
	data, err := yaml.Marshal(cache)
	if err == nil {
		// write privately owned by the user
		err = os.WriteFile(filename, data, 0600)
	}
	return err
}

// Retrieve() implements the Provider interface, returning the cached credential if is not expired,
// otherwise fetching the credential from the underlying Provider and caching the results on disk
// with an expiration time.
func (f *FileCacheProvider) Retrieve() (credentials.Value, error) {
	if !f.cachedCredential.IsExpired() {
		// use the cached credential
		return f.cachedCredential.Credential, nil
	}
	fmt.Fprintf(os.Stderr, "No cached credential available.  Refreshing...\n")
	// fetch the credentials from the underlying Provider
	credential, err := f.credentials.Get()
	if err != nil {
		return credential, err
	}
	expiration, err := f.credentials.ExpiresAt()
	if err != nil {
		// credential doesn't support expiration time, so can't cache, but still return the credential
		fmt.Fprintf(os.Stderr, "Unable to cache credential: %v\n", err)
		return credential, nil
	}
	// underlying provider supports Expirer interface, so we can cache
	filename, err := cacheFilename()
	if err != nil {
		return credential, err
	}
	f.cachedCredential = cachedCredential{
		Credential:  credential,
		Expiration:  expiration,
		currentTime: nil,
	}
	// overwrite whatever was there before. we don't care about multiple creds for various clusters.
	// if user switches to another role and another profile they have to re-authenticate.
	cache, _ := readCacheFile(filename)
	cache.Put(f.profile, f.cachedCredential)
	if err := writeCache(filename, cache); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update credential cache %s: %v\n", filename, err)
	}
	fmt.Fprintf(os.Stderr, "Updated cached credential\n")
	return credential, err
}

// IsExpired() implements the Provider interface, deferring to the cached credential first,
// but fall back to the underlying Provider if it is expired.
func (f *FileCacheProvider) IsExpired() bool {
	return f.cachedCredential.IsExpired() && f.credentials.IsExpired()
}

// ExpiresAt implements the Expirer interface, and gives access to the expiration time of the credential
func (f *FileCacheProvider) ExpiresAt() time.Time {
	return f.cachedCredential.Expiration
}

func cacheFilename() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filename := filepath.Join(home, ".eksctl", "cache", "credentials.yaml")
	return filename, nil
}
