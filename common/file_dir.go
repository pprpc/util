package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PathIsExist check path
func PathIsExist(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		return false
	}
	return true
}

// IsDir .
func IsDir(p string) bool {
	fi, err := os.Lstat(p)
	if err != nil {
		return false
	}
	md := fi.Mode()
	if md.IsDir() {
		return true
	}
	return false
}

// IsFile .
func IsFile(p string) bool {
	fi, err := os.Lstat(p)
	if err != nil {
		return false
	}
	md := fi.Mode()
	if md.IsDir() == false {
		return true
	}
	return false
}

// CreateDir create dir
func CreateDir(dirPath string) (err error) {
	cdir, _ := filepath.Split(dirPath)
	if 0 != len(cdir) {
		err = os.MkdirAll(cdir, 0755)
		return
	}
	return
}

func CreateDirv2(dirPath string) (err error) {
	err = os.MkdirAll(dirPath, 0755)
	return
}

// CreateFile create or update file
func CreateFile(fp string, body []byte) (err error) {
	err = CreateDir(fp)
	if err != nil {
		err = fmt.Errorf("CreateDir(), %s", err)
		return
	}
	var fd *os.File
	fd, err = os.Create(fp)
	if err != nil {
		err = fmt.Errorf("os.Create(), %s", err)
		return
	}
	defer fd.Close()

	_, err = fd.Write(body)

	return
}

// FileMd5 获取文件MD5值
func FileMd5(fileName string) (sMd5 string, err error) {
	f, lerr := os.Open(fileName)
	if lerr != nil {
		err = lerr
		return
	}
	defer f.Close()

	h := md5.New()
	if _, lerr := io.Copy(h, f); lerr != nil {
		err = lerr
		return
	}
	sMd5 = fmt.Sprintf("%x", h.Sum(nil))
	return
}
