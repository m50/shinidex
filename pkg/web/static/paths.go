package static

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/gookit/slog"
)

const scriptFSPath = "assets/scripts.js"
const styleFSPath = "assets/style.css"

var (
	versionedScriptPath = ""
	versionedStylePath  = ""
	inInit              = true
)

func init() {
	// Call the get methods on init, so that we can cache the versioned names
	GetScriptPath()
	GetStylePath()
	inInit = false
}

func GetScriptPath() string {
	if versionedScriptPath != "" {
		return versionedScriptPath
	}
	h, err := getFileHash(scriptFSPath)
	if err != nil {
		if !inInit {
			slog.Warnf("unable to generate versioned script path: %v", err)
		}
		return "/assets/scripts.js"
	}
	versionedScriptPath = fmt.Sprintf("/assets/scripts.%s.js", h)
	return versionedScriptPath
}

func GetStylePath() string {
	if versionedStylePath != "" {
		return versionedStylePath
	}
	h, err := getFileHash(styleFSPath)
	if err != nil {
		if !inInit {
			slog.Warnf("unable to generate versioned style path: %v", err)
		}
		return "/assets/style.css"
	}
	versionedStylePath = fmt.Sprintf("/assets/style.%s.css", h)
	return versionedStylePath
}

func getFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil))[:10], nil
}
