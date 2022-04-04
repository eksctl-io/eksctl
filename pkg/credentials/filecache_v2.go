package credentials

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/kris-nova/logger"
	"github.com/spf13/afero"
)

// FileCacheV2 is a file-based credentials cache for AWS credentials that can expire,
// satisfying the aws.CredentialsProvider interface.
// It is meant to be wrapped with aws.CredentialsCache. The cache is per profile.
type FileCacheV2 struct {
	provider      aws.CredentialsProvider
	profileName   string
	cacheFilePath string
	fs            afero.Fs
	newFlock      FlockFunc
	clock         Clock

	creds *aws.Credentials
	mu    sync.Mutex
}

// NewFileCacheV2 initializes the cache and returns a *FileCacheV2.
func NewFileCacheV2(provider aws.CredentialsProvider, profileName string, fs afero.Fs, newFlock FlockFunc, clock Clock, cacheFilePath string) (*FileCacheV2, error) {
	if err := initializeCache(fs, cacheFilePath); err != nil {
		return nil, fmt.Errorf("error initializing credentials cache: %w", err)
	}
	return &FileCacheV2{
		provider:      provider,
		profileName:   profileName,
		cacheFilePath: cacheFilePath,
		fs:            fs,
		newFlock:      newFlock,
		clock:         clock,
	}, nil
}

func toAWSCredentials(c cachedCredential) *aws.Credentials {
	return &aws.Credentials{
		AccessKeyID:     c.Credential.AccessKeyID,
		SecretAccessKey: c.Credential.SecretAccessKey,
		SessionToken:    c.Credential.SessionToken,
		Source:          c.Credential.ProviderName,
		CanExpire:       !c.Expiration.IsZero(),
		Expires:         c.Expiration,
	}
}

// Retrieve implements aws.CredentialsProvider.
func (f *FileCacheV2) Retrieve(ctx context.Context) (aws.Credentials, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.creds == nil {
		cacheFile, err := readCacheFile(f.fs, f.cacheFilePath, f.newFlock)
		if err != nil {
			logger.Warning("error reading credentials cache: %v", err)
		} else {
			creds, ok := cacheFile.ProfileMap[f.profileName]
			if ok {
				f.creds = toAWSCredentials(creds)
			}
		}
	}

	if f.creds != nil && f.creds.CanExpire && f.creds.Expires.After(f.clock.Now().Round(0)) {
		return *f.creds, nil
	}

	creds, err := f.provider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}
	f.creds = &creds

	if !creds.CanExpire {
		return creds, nil
	}

	cache, err := readCacheFile(f.fs, f.cacheFilePath, f.newFlock)
	if err != nil {
		logger.Warning("error reading cache file: %v", err)
		return creds, nil
	}
	cache.Put(f.profileName, cachedCredential{
		Credential: credentials.Value{
			AccessKeyID:     creds.AccessKeyID,
			SecretAccessKey: creds.SecretAccessKey,
			SessionToken:    creds.SessionToken,
			ProviderName:    creds.Source,
		},
		Expiration: creds.Expires,
	})

	if err := writeCache(f.fs, f.cacheFilePath, f.newFlock, cache); err != nil {
		logger.Warning("failed to update credentials cache: %v", err)
	}

	return creds, nil
}
