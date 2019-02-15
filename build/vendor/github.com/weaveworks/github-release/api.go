package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	API_URL = "https://api.github.com"
	GH_URL  = "https://github.com"
)

// materializeFile takes a physical file or stream (named pipe, user input,
// ...) and returns an io.Reader and the number of bytes that can be read
// from it.
func materializeFile(f *os.File) (io.Reader, int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}

	// If the file is actually a char device (like user typed input)
	// or a named pipe (like a streamed in file), buffer it up.
	//
	// When uploading a file, you need to either explicitly set the
	// Content-Length header or send a chunked request. Since the
	// github upload server doesn't accept chunked encoding, we have
	// to set the size of the file manually. Since a stream doesn't have a
	// predefined length, it's read entirely into a byte buffer.
	if fi.Mode()&(os.ModeCharDevice|os.ModeNamedPipe) == 1 {
		vprintln("input was a stream, buffering up")

		var buf bytes.Buffer
		n, err := buf.ReadFrom(f)
		if err != nil {
			return nil, 0, errors.New("req: could not buffer up input stream: " + err.Error())
		}
		return &buf, n, err
	}

	// We know the os.File is most likely an actual file now.
	n, err := GetFileSize(f)
	return f, n, err
}

/* create a new request that sends the auth token */
func NewAuthRequest(method, url, bodyType, token string, headers map[string]string, body io.Reader) (*http.Request, error) {
	vprintln("creating request:", method, url, bodyType, token)

	var n int64 // content length
	var err error
	if f, ok := body.(*os.File); ok {
		// Retrieve the content-length and buffer up if necessary.
		body, n, err = materializeFile(f)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if n != 0 {
		vprintln("setting content-length to", n)
		req.ContentLength = n
	}

	if bodyType != "" {
		req.Header.Set("Content-Type", bodyType)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func DoAuthRequest(method, url, bodyType, token string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := NewAuthRequest(method, url, bodyType, token, headers, body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GithubGet(uri string, v interface{}) error {
	resp, err := http.Get(ApiURL() + uri)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("could not fetch releases, %v", err)
	}

	vprintln("GET", ApiURL()+uri, "->", resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github did not response with 200 OK but with %v", resp.Status)
	}

	if VERBOSITY == 0 {
		if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("could not unmarshall JSON into Release struct, %v", err)
		}
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		vprintln("BODY", string(body))

		if err = json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("could not unmarshall JSON into Release struct, %v", err)
		}
	}

	return nil
}

func ApiURL() string {
	if "" == EnvApiEndpoint {
		return API_URL
	} else {
		return EnvApiEndpoint
	}
}
