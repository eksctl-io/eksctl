package filecache

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

type stubProvider struct {
	creds   credentials.Value
	expired bool
	err     error
}

func (s *stubProvider) Retrieve() (credentials.Value, error) {
	s.expired = false
	s.creds.ProviderName = "stubProvider"
	return s.creds, s.err
}

func (s *stubProvider) IsExpired() bool {
	return s.expired
}

type stubProviderExpirer struct {
	stubProvider
	expiration time.Time
}

func (s *stubProviderExpirer) ExpiresAt() time.Time {
	return s.expiration
}

type testFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fs *testFileInfo) Name() string       { return fs.name }
func (fs *testFileInfo) Size() int64        { return fs.size }
func (fs *testFileInfo) Mode() os.FileMode  { return fs.mode }
func (fs *testFileInfo) ModTime() time.Time { return fs.modTime }
func (fs *testFileInfo) IsDir() bool        { return fs.Mode().IsDir() }
func (fs *testFileInfo) Sys() interface{}   { return nil }

type testFS struct {
	filename string
	fileinfo testFileInfo
	data     []byte
	err      error
	perm     os.FileMode
}

func (t *testFS) Stat(filename string) (os.FileInfo, error) {
	t.filename = filename
	return &t.fileinfo, t.err
}

func (t *testFS) ReadFile(filename string) ([]byte, error) {
	t.filename = filename
	return t.data, t.err
}

func (t *testFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	t.filename = filename
	t.data = data
	t.perm = perm
	return t.err
}

func (t *testFS) MkdirAll(path string, perm os.FileMode) error {
	t.filename = path
	t.perm = perm
	return t.err
}

func (t *testFS) reset() {
	t.filename = ""
	t.fileinfo = testFileInfo{}
	t.data = []byte{}
	t.err = nil
	t.perm = 0600
}

type testEnv struct {
	values map[string]string
}

func (e *testEnv) Getenv(key string) string {
	return e.values[key]
}

func (e *testEnv) LookupEnv(key string) (string, bool) {
	value, ok := e.values[key]
	return value, ok
}

func (e *testEnv) reset() {
	e.values = map[string]string{}
}

type testFilelock struct {
	ctx        context.Context
	retryDelay time.Duration
	success    bool
	err        error
}

func (l *testFilelock) Unlock() error {
	return nil
}

func (l *testFilelock) TryLockContext(ctx context.Context, retryDelay time.Duration) (bool, error) {
	l.ctx = ctx
	l.retryDelay = retryDelay
	return l.success, l.err
}

func (l *testFilelock) TryRLockContext(ctx context.Context, retryDelay time.Duration) (bool, error) {
	l.ctx = ctx
	l.retryDelay = retryDelay
	return l.success, l.err
}

func (l *testFilelock) reset() {
	l.ctx = context.TODO()
	l.retryDelay = 0
	l.success = true
	l.err = nil
}

func getMocks() (tf *testFS, te *testEnv, testFlock *testFilelock) {
	tf = &testFS{}
	tf.reset()
	f = tf
	te = &testEnv{}
	te.reset()
	e = te
	testFlock = &testFilelock{}
	testFlock.reset()
	newFlock = func(filename string) filelock {
		return testFlock
	}
	return
}

func makeCredential() credentials.Value {
	return credentials.Value{
		AccessKeyID:     "AKID",
		SecretAccessKey: "SECRET",
		SessionToken:    "TOKEN",
		ProviderName:    "stubProvider",
	}
}

func validateFileCacheProvider(t *testing.T, p Provider, err error, c *credentials.Credentials) {
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.credentials != c {
		t.Errorf("Credentials not copied")
	}
	if p.cacheKey.clusterID != "CLUSTER" {
		t.Errorf("clusterID not copied")
	}
	if p.cacheKey.profile != "PROFILE" {
		t.Errorf("profile not copied")
	}
	if p.cacheKey.roleARN != "ARN" {
		t.Errorf("roleARN not copied")
	}
}

func TestCacheFilename(t *testing.T) {
	_, te, _ := getMocks()

	te.values["HOME"] = "homedir"        // unix
	te.values["USERPROFILE"] = "homedir" // windows

	filename := CacheFilename()
	expected := "homedir/.eksctl/cache/credentials.yaml"
	if filename != expected {
		t.Errorf("Incorrect default cacheFilename, expected %s, got %s",
			expected, filename)
	}

	te.values["EKSCTL_CACHE_FILE"] = "special.yaml"
	filename = CacheFilename()
	expected = "special.yaml"
	if filename != expected {
		t.Errorf("Incorrect custom cacheFilename, expected %s, got %s",
			expected, filename)
	}
}

func TestNewFileCacheProvider_Missing(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	tf, _, _ := getMocks()

	// missing cache file
	tf.err = os.ErrNotExist
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)
	if !p.cachedCredential.IsExpired() {
		t.Errorf("missing cache file should result in expired cached credential")
	}
	tf.err = nil
}

func TestNewFileCacheProvider_BadPermissions(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	tf, _, _ := getMocks()

	// bad permissions
	tf.fileinfo.mode = 0777
	_, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	if err == nil {
		t.Errorf("Expected error due to public permissions")
	}
	if tf.filename != CacheFilename() {
		t.Errorf("unexpected file checked, expected %s, got %s",
			CacheFilename(), tf.filename)
	}
}

func TestNewFileCacheProvider_Unlockable(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	_, _, testFlock := getMocks()

	// unable to lock
	testFlock.success = false
	testFlock.err = errors.New("lock stuck, needs wd-40")
	_, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	if err == nil {
		t.Errorf("Expected error due to lock failure")
	}
	testFlock.success = true
	testFlock.err = nil
}

func TestNewFileCacheProvider_Unreadable(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	tf, _, _ := getMocks()

	// unable to read existing cache
	tf.err = errors.New("read failure")
	_, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	if err == nil {
		t.Errorf("Expected error due to read failure")
	}
	tf.err = nil
}

func TestNewFileCacheProvider_Unparseable(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	tf, _, _ := getMocks()

	// unable to parse yaml
	tf.data = []byte("invalid: yaml: file")
	_, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	if err == nil {
		t.Errorf("Expected error due to bad yaml")
	}
}

func TestNewFileCacheProvider_Empty(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	_, _, _ = getMocks()

	// successfully parse existing but empty cache file
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)
	if !p.cachedCredential.IsExpired() {
		t.Errorf("empty cache file should result in expired cached credential")
	}
}

func TestNewFileCacheProvider_ExistingCluster(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	tf, _, _ := getMocks()

	// successfully parse existing cluster without matching arn
	tf.data = []byte(`clusters:
  CLUSTER:
`)
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)
	if !p.cachedCredential.IsExpired() {
		t.Errorf("missing arn in cache file should result in expired cached credential")
	}
}

func TestNewFileCacheProvider_ExistingARN(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	tf, _, _ := getMocks()

	// successfully parse cluster with matching arn
	tf.data = []byte(`clusters:
  CLUSTER:
    PROFILE:
      ARN:
        credential:
          accesskeyid: ABC
          secretaccesskey: DEF
          sessiontoken: GHI
          providername: JKL
        expiration: 2018-01-02T03:04:56.789Z
`)
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)
	if p.cachedCredential.Credential.AccessKeyID != "ABC" || p.cachedCredential.Credential.SecretAccessKey != "DEF" ||
		p.cachedCredential.Credential.SessionToken != "GHI" || p.cachedCredential.Credential.ProviderName != "JKL" {
		t.Errorf("cached credential not extracted correctly")
	}
	// fiddle with clock
	p.cachedCredential.currentTime = func() time.Time {
		return time.Date(2017, 12, 25, 12, 23, 45, 678, time.UTC)
	}
	if p.cachedCredential.IsExpired() {
		t.Errorf("Cached credential should not be expired")
	}
	if p.IsExpired() {
		t.Errorf("Cache credential should not be expired")
	}
	expectedExpiration := time.Date(2018, 01, 02, 03, 04, 56, 789000000, time.UTC)
	if p.ExpiresAt() != expectedExpiration {
		t.Errorf("Credential expiration time is not correct, expected %v, got %v",
			expectedExpiration, p.ExpiresAt())
	}
}

func TestFileCacheProvider_Retrieve_NoExpirer(t *testing.T) {
	providerCredential := makeCredential()
	c := credentials.NewCredentials(&stubProvider{
		creds: providerCredential,
	})

	tf, _, _ := getMocks()

	// initialize from missing cache file
	tf.err = os.ErrNotExist
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)

	credential, err := p.Retrieve()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if credential != providerCredential {
		t.Errorf("Cache did not return provider credential, got %v, expected %v",
			credential, providerCredential)
	}
}

func makeExpirerCredentials() (providerCredential credentials.Value, expiration time.Time, c *credentials.Credentials) {
	providerCredential = makeCredential()
	expiration = time.Date(2020, 9, 19, 13, 14, 0, 1000000, time.UTC)
	c = credentials.NewCredentials(&stubProviderExpirer{
		stubProvider{
			creds: providerCredential,
		},
		expiration,
	})
	return
}

func TestFileCacheProvider_Retrieve_WithExpirer_Unlockable(t *testing.T) {
	providerCredential, _, c := makeExpirerCredentials()

	tf, _, testFlock := getMocks()

	// initialize from missing cache file
	tf.err = os.ErrNotExist
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)

	// retrieve credential, which will fetch from underlying Provider
	// fail to get write lock
	testFlock.success = false
	testFlock.err = errors.New("lock stuck, needs wd-40")
	credential, err := p.Retrieve()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if credential != providerCredential {
		t.Errorf("Cache did not return provider credential, got %v, expected %v",
			credential, providerCredential)
	}
}

func TestFileCacheProvider_Retrieve_WithExpirer_Unwritable(t *testing.T) {
	providerCredential, expiration, c := makeExpirerCredentials()

	tf, _, _ := getMocks()

	// initialize from missing cache file
	tf.err = os.ErrNotExist
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)

	// retrieve credential, which will fetch from underlying Provider
	// fail to write cache
	tf.err = errors.New("can't write cache")
	credential, err := p.Retrieve()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if credential != providerCredential {
		t.Errorf("Cache did not return provider credential, got %v, expected %v",
			credential, providerCredential)
	}
	if tf.filename != CacheFilename() {
		t.Errorf("Wrote to wrong file, expected %v, got %v",
			CacheFilename(), tf.filename)
	}
	if tf.perm != 0600 {
		t.Errorf("Wrote with wrong permissions, expected %o, got %o",
			0600, tf.perm)
	}
	expectedData := []byte(`clusters:
  CLUSTER:
    PROFILE:
      ARN:
        credential:
          accesskeyid: AKID
          secretaccesskey: SECRET
          sessiontoken: TOKEN
          providername: stubProvider
        expiration: ` + expiration.Format(time.RFC3339Nano) + `
`)
	if !bytes.Equal(tf.data, expectedData) {
		t.Errorf("Wrong data written to cache, expected: %s, got %s",
			expectedData, tf.data)
	}
}

func TestFileCacheProvider_Retrieve_WithExpirer_Writable(t *testing.T) {
	providerCredential, _, c := makeExpirerCredentials()

	tf, _, _ := getMocks()

	// initialize from missing cache file
	tf.err = os.ErrNotExist
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)
	tf.err = nil

	// retrieve credential, which will fetch from underlying Provider
	// same as TestFileCacheProvider_Retrieve_WithExpirer_Unwritable,
	// but write to disk (code coverage)
	credential, err := p.Retrieve()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if credential != providerCredential {
		t.Errorf("Cache did not return provider credential, got %v, expected %v",
			credential, providerCredential)
	}
}

func TestFileCacheProvider_Retrieve_CacheHit(t *testing.T) {
	c := credentials.NewCredentials(&stubProvider{})

	tf, _, _ := getMocks()

	// successfully parse cluster with matching arn
	tf.data = []byte(`clusters:
  CLUSTER:
    PROFILE:
      ARN:
        credential:
          accesskeyid: ABC
          secretaccesskey: DEF
          sessiontoken: GHI
          providername: JKL
        expiration: 2018-01-02T03:04:56.789Z
`)
	p, err := NewFileCacheProvider("CLUSTER", "PROFILE", "ARN", c)
	validateFileCacheProvider(t, p, err, c)

	// fiddle with clock
	p.cachedCredential.currentTime = func() time.Time {
		return time.Date(2017, 12, 25, 12, 23, 45, 678, time.UTC)
	}

	credential, err := p.Retrieve()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if credential.AccessKeyID != "ABC" || credential.SecretAccessKey != "DEF" ||
		credential.SessionToken != "GHI" || credential.ProviderName != "JKL" {
		t.Errorf("cached credential not returned")
	}
}
