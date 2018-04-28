package clccam

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"

	"github.com/pkg/errors"
)

const (
	// Name of the file to store the last bearer-token
	tokenFile = "cam.token"
)

// Token is the CAM JWT Authorization token
type Token string

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
