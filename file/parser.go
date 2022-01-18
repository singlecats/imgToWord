package file

import (
	"imageToWord/zip"
	"os"
	"path/filepath"
)

func ParseFile(path string, info os.FileInfo, err error) error {
	if IsDir(path) {
		return nil
	}
	return nil
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func IsImage(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".png":
		return true
	case ".jpg":
		return true
	case ".jpeg":
		return true
	default:
		return false
	}
}
func CheckDir(path string)  {
	if ok, _ := zip.PathExists(path); !ok {
		os.MkdirAll(path, 0777)
	}
}