package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func DirExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil // Path exists, return true if it's a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func DirExistsAndValidPerms(path string, Perm string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if checkPermitionOther(info.Mode(), Perm) == false {
			return false, errors.New(fmt.Sprint("[util] Permition denied : File ", path, " does not have permition ", Perm))
		}
		return info.IsDir(), nil // Path exists, return true if it's a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func checkPermitionOther(m fs.FileMode, reqPermStr string) bool {
	UnixStylePermStr := m.String()
	PermOtherUser := UnixStylePermStr[len(UnixStylePermStr)-3:]

	for _, per := range reqPermStr {
		if strings.ContainsRune(PermOtherUser, per) == false {
			return false
		}
	}
	return true
}

func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil // Path exists, return true if it's a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func FileExistsAndValidPerms(path string, Perm string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if checkPermitionOther(info.Mode(), Perm) == false {
			return false, errors.New(fmt.Sprint("[util] Permition Denied File", path, "does not have permition", Perm))
		}
		return !info.IsDir(), nil // Path exists, return true if it's a directory
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil // Path does not exist
	}
	return false, err // Other error (e.g., permission denied)
}

func SaveFileFromBuf(path string, src io.Reader) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return errors.New("[exec error] failed to open stdout file")
	}
	defer f.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllFilesInDir(dir_path string) error {
	d, err := os.Open(dir_path)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir_path, name))
        if err != nil {
            return err
        }
    }
    return nil
}
