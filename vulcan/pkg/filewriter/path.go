package filewriter

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

var (
	destDir string
	tmpDir  string
)

func Init(dest, sep string) {
	destDir = dest
	tmpDir = sep
}

func WriteFile(filename string, content []byte) error {
	if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
		return errors.Wrapf(err, "%s already exists", filename)
	}

	subDir := strings.ReplaceAll(filename, filepath.Base(filename), "")
	if !PathExists(subDir) {
		if err := os.MkdirAll(subDir, 0750); err != nil {
			return errors.Wrapf(err, "create directory failed. dir: %s file: %s", subDir, filename)
		}
	}

	if err := os.WriteFile(filename, content, 0600); err != nil {
		return errors.Wrapf(err, "write file failed. file: %s", filename)
	}

	return nil
}

func RebuildDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return errors.Wrapf(err, "remove directory failed. dir: %s", dir)
	}

	if err := os.MkdirAll(dir, 0750); err != nil {
		return errors.Wrapf(err, "create directory failed. dir: %s", dir)
	}

	return nil
}

func CreateDir(path string) (existed bool, err error) {
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		return true, nil
	}

	if err := os.MkdirAll(path, 0750); err != nil {
		return false, errors.Wrapf(err, "create directory failed. path: %s", path)
	}

	return false, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !os.IsNotExist(err)
}

func BasePath() string {
	baseDir := os.Getenv("ROMA_ROOT")
	if len(baseDir) == 0 {
		baseDir, _ = os.Getwd()
	}

	return filepath.FromSlash(filepath.Clean(baseDir))
}

func GetCurrentAbPath() (string, error) {
	baseDir, err := getCurrentAbPathByExecutable()
	if err != nil {
		return "", err
	}

	tmpDir, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		return "", errors.Wrapf(err, "failed to get temp dir")
	}

	if strings.Contains(baseDir, tmpDir) {
		return getCurrentAbPathByCaller(), nil
	}

	return filepath.FromSlash(filepath.Clean(baseDir)), nil
}

// get current executable file absolute path
func getCurrentAbPathByExecutable() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get executable path")
	}

	res, err := filepath.EvalSymlinks(filepath.Dir(exePath))
	if err != nil {
		return "", errors.Wrapf(err, "failed to get executable path")
	}

	return res, nil
}

// get current executable file absolute path on executing 'go run'
func getCurrentAbPathByCaller() string {
	var abPath string

	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = filepath.Dir(filename)
	}

	return abPath
}

func SprintGenPath(path string) string {
	if tmpDir == "" {
		return path
	}

	parts := strings.Split(path, filepath.Clean(tmpDir))
	if len(parts) == 2 {
		return filepath.Join(destDir, parts[1])
	}

	return path
}
