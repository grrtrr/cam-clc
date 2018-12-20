package clccam

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/grrtrr/clccam/logger"
	"github.com/pkg/errors"
)

// Name of the file to cache the last-used bearer-token.
const tokenFile = "cam.token"

// Token is the CAM JWT Authorization token
type Token string

func (t Token) String() string {
	c, err := t.Claims()
	if err != nil {
		return fmt.Sprintf("invalid CAM token (%s)", err)
	}
	return c.String()
}

// NewClient returns a CAM client that uses @t as authorization/bearer token.
func (t Token) NewClient(options ...ClientOption) *Client {
	return NewClient(RequestOptions(Headers(map[string]string{
		"Authorization": "Bearer " + string(t),
	}))).With(options...)
}

// Decode attempts to parse @t, returning an error if it fails to parse.
func (t Token) Decode() (payload []byte, err error) {
	payload, _, err = jose.DecodeBytes(string(t), camJwtTokenPublicKey())
	return payload, err
}

// Claims extracts the CAM claims payload from @t.
func (t Token) Claims() (*Claims, error) {
	var pl Claims

	if payload, err := t.Decode(); err != nil {
		return nil, err
	} else if err := json.Unmarshal(payload, &pl); err != nil {
		return nil, errors.Wrapf(err, "unable to extract CAM token claims payload")
	}
	return &pl, nil
}

// LoadToken attempts to load a CAM token from the environment variable $CAM_TOKEN or $tokenFile.
func LoadToken() (Token, error) {
	var tokenPath = path.Join(GetClcHome(), tokenFile)

	if token := os.Getenv("CAM_TOKEN"); token != "" {
		return Token(token), nil
	} else if _, err := os.Stat(tokenPath); err == nil {
		fd, err := os.Open(tokenPath)
		if err != nil {
			return "", errors.Errorf("failed to load token from %s: %s", tokenPath, err)
		}
		defer fd.Close()

		content, err := ioutil.ReadAll(fd)
		if err != nil {
			return "", errors.Errorf("failed to read %s: %s", tokenPath, err)
		}
		return Token(bytes.TrimSpace(content)), nil
	}
	return "", errors.Errorf("no valid token configuration found in %s", GetClcHome())
}

// SaveToken saves @token to file.
func (t Token) Save() error {
	return writeCLCdata(tokenFile, []byte(t), 0600)
}

// writeCLCitem writes @data to CLC_HOME/fileName
func writeCLCdata(fileName string, data []byte, perm os.FileMode) error {
	var clcHome = GetClcHome()

	if _, err := os.Stat(clcHome); os.IsNotExist(err) {
		if err = os.MkdirAll(clcHome, 0700); err != nil {
			return errors.Errorf("failed to create CLC directory %s: %s", clcHome, err)
		}
	}
	return ioutil.WriteFile(path.Join(clcHome, fileName), data, perm)
}

// GetClcHome returns the path to the CLC CAM configuration directory, which is the same
// as used by, and compatible with, clc-go-cli (including the CLC_HOME environment variable).
func GetClcHome() string {
	if clcHome := os.Getenv("CLC_HOME"); clcHome != "" {
		return clcHome
	}

	u, err := user.Current()
	if err != nil {
		log.Fatalf("failed to look up current user: %s", err)
	}

	if runtime.GOOS == "windows" {
		return path.Join(u.HomeDir, "clc")
	} else {
		return path.Join(u.HomeDir, ".clc")
	}
}

// camJwtTokenPublicKey returns the RSA Public Key that can validate the signatures of CAM-generated access tokens.
func camJwtTokenPublicKey() *rsa.PublicKey {
	const camJwtPubKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvfXAptp4XtBpIlXPzu0i
Y7trJ5XOlgFpIw742q56AMXi1s9M1KS3qbZwz1Bkk7UX3SS+ZdyXvb1M23jzu7Ji
lenUBBEea974eNm3mIdwTcuVeuVf3Xn7plU59eJNTzMCgz/OV9Zo6YNsHpHnBGVE
mBfstcNCuufbNC80zzE1YEthkIsPcoJgl4imUH6nl3sHx8ndMsz4MBnLkHsz0pXG
53bmwKJF7kh/gYL/5+WJmzwsh1tsGWKkDr1pPedW0oNJLADy3MfmA/kFaa7NRL0z
p7w9pVV/CO5J6XrtVoaVJz1A31pAc85qez8qZluGJ9SqZhM2XgmBiaDEvYSOvCED
7QIDAQAB
-----END PUBLIC KEY-----
`
	block, rem := pem.Decode([]byte(camJwtPubKey))
	if block == nil {
		logger.Fatalf("unable to extract CAM JWT public key from %q", camJwtPubKey)
	} else if len(rem) > 0 {
		logger.Fatalf("extra data at end of CAM JWT public key: %q", string(rem))
	} else if block.Type != "PUBLIC KEY" {
		logger.Fatalf("CAM JWT public key has inconsistent type %q", block.Type)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Fatalf("failed to parse CAM JWT public key: %s", err)
	}
	return pub.(*rsa.PublicKey)
}
