package zip

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var tmpPath = filepath.Join(".", "tmp")

func ParseZip(zipFile []string) []string {
	var resArr []string
	resChannel := make([]chan string, len(zipFile))
	for i, file := range zipFile {
		resChannel[i] = make(chan string)
		go func(ch chan string, file string) {
			defer func() {
				if e := recover(); e != nil {
					fmt.Println(e)
					ch <- "-1"
				}
			}()
			if FileExist(file) {
				currentTime := time.Now()
				distinctPath := filepath.Join(tmpPath, currentTime.Format("2006-01-02"), getFileName(file) + "_")
				CheckDir(distinctPath)
				log.Println(distinctPath)
				err := DeCompress(file, distinctPath)
				if err != nil {
					fmt.Println(err)
					//log.Fatal("文件解压失败 %s", err)
					panic(err)
				}
				ch <- distinctPath
			}
		}(resChannel[i], file)
	}
	for _, ch := range resChannel {
		resArr = append(resArr, <-ch)
	}
	fmt.Println(resArr)
	return resArr
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func getFileName(file string) string {
	filenameall := path.Base(file)
	filesuffix := path.Ext(file)
	return filenameall[0 : len(filenameall)-len(filesuffix)]
}

func CheckDir(path string) {
	if ok, _ := PathExists(path); !ok {
		os.MkdirAll(path, 0777)
	}
}
