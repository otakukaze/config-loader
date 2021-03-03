package confloader

import (
	"os"
	"path"
	"runtime"
	"strings"
)

// ParsePath - parse file path to absPath
func ParsePath(dst string) string {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	if []rune(dst)[0] == '~' {
		home := UserHomeDir()
		if len(home) > 0 {
			dst = strings.Replace(dst, "~", home, -1)
		}
	}

	if path.IsAbs(dst) {
		dst = path.Clean(dst)
		return dst
	}

	str := path.Join(wd, dst)
	str = path.Clean(str)
	return str
}

// UserHomeDir - get user home directory
func UserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

// CheckExists - check file exists
func CheckExists(filePath string, allowDir bool) bool {
	filePath = ParsePath(filePath)
	stat, err := os.Stat(filePath)
	if err != nil && !os.IsExist(err) {
		return false
	}
	if !allowDir && stat.IsDir() {
		return false
	}
	return true
}

// IsDir -
func IsDir(filePath string) bool {
	filePath = ParsePath(filePath)
	stat, err := os.Stat(filePath)
	if err != nil && !os.IsExist(err) {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	return true
}
