package filewriter

import (
	"log"
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

	subDir := strings.Replace(filename, filepath.Base(filename), "", -1)
	if !PathExists(subDir) {
		if err := os.MkdirAll(subDir, 0755); err != nil {
			return errors.Wrapf(err, "create directory failed. dir: %s file: %s", subDir, filename)
		}
	}

	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		return errors.Wrapf(err, "write file failed. file: %s", filename)
	}
	return nil
}

func RebuildDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return errors.Wrapf(err, "remove directory failed. dir: %s", dir)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.Wrapf(err, "create directory failed. dir: %s", dir)
	}
	return nil
}

func CreateDir(path string) (existed bool, err error) {
	if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
		return true, nil
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return false, errors.Wrapf(err, "create directory failed. path: %s", path)
	}
	return false, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func BasePath() string {
	baseDir := os.Getenv("ROMA_ROOT")
	if len(baseDir) == 0 {
		baseDir, _ = os.Getwd()
	}
	return filepath.FromSlash(filepath.Clean(baseDir))
}

func GetCurrentAbPath() string {
	baseDir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(baseDir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	baseDir = filepath.Clean(baseDir)
	return filepath.FromSlash(baseDir)
}

// get current executable file absolute path
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
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
