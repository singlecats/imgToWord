package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imageToWord/doc"
	"imageToWord/file"
	"imageToWord/zip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Msg struct {
	Code    string
	Message string
}

func init()  {
	fmt.Println("服务启动成功....")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	// 写入日志的文件
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.Static("/assets", "./public")
	r.LoadHTMLGlob("./public/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/upload", func(c *gin.Context) {
		defer func(c *gin.Context) {
			ret := Msg{
				Code:    "400",
				Message: "",
			}

			err := recover()
			switch err.(type) {
			case runtime.Error: // 运行时错误
				ret.Message = "发生错误了"
				c.JSON(http.StatusOK, ret)
			default: // 非运行时错误

			}
		}(c)
		form, _ := c.MultipartForm()
		files := form.File["file"]
		currentTime := time.Now()
		tmpFile := filepath.Join(".", "tmp", currentTime.Format("2006-01-02"))
		file.CheckDir(tmpFile)
		var zipArr []string
		for _, dirFile := range files {
			log.Println(dirFile.Filename)
			dst := filepath.Join(tmpFile, dirFile.Filename)
			// 上传文件到指定的 dst.
			err := c.SaveUploadedFile(dirFile, dst)
			if err != nil {
				continue
			}
			log.Println(dst)
			zipArr = append(zipArr, dst)
		}
		ok, msg := start(zipArr)
		if ok {
			ret := Msg{
				Code:    "200",
				Message: fmt.Sprintf("%d 个文件处理成功", len(files)),
			}
			c.JSON(http.StatusOK, ret)
		} else {
			ret := Msg{
				Code:    "400",
				Message: msg,
			}
			c.JSON(http.StatusOK, ret)
		}
	})
	startCheckTmp()
	r.Run(":9090") // listen and serve on 0.0.0.0:8080

}

var wg sync.WaitGroup

func start(zipArr []string) (bool, string) {
	resArr := zip.ParseZip(zipArr)

	if len(resArr) == 0 {
		fmt.Printf("无有效文件")
		return false, "无有效文件"
	}
	total := 0

	for _, tmpDir := range resArr {
		if ok, _ := zip.PathExists(tmpDir); ok {
			wg.Add(1)
			total++
		}
	}
	if total == 0 {
		return false, "无有效文件"
	}
	errMsg := make([]chan string, total)
	fmt.Println(len(errMsg))
	for i, tmpDir := range resArr {
		if ok, _ := zip.PathExists(tmpDir); ok {
			errMsg[i] = make(chan string)
			go process(tmpDir, errMsg[i])
		}
	}

	errMsgStr := ""
	//fmt.Println(<-errMsg[0])
	for _, msg := range errMsg {
		select {
		case message := <-msg:

			errMsgStr += message
			break
		}
	}
	if errMsgStr != "" {
		return false, errMsgStr
	}
	wg.Wait()
	return true, ""
}

func process(tmpDir string, errMsg chan string) {
	defer wg.Done()
	defer close(errMsg)
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			errMsg <- "生成doc文档失败"
		} else {
			errMsg <- ""
		}
	}()
	if ok, _ := zip.PathExists(tmpDir); !ok {
		return
	}
	_, dirName := filepath.Split(tmpDir)
	dirArr := strings.Split(dirName, "_")
	currentTime := time.Now()
	downloadFile := filepath.Join(".", "download", currentTime.Format("2006-01-02"))
	zip.CheckDir(downloadFile)
	saveName := filepath.Join(downloadFile, currentTime.Format("15_4_5_")+dirArr[0]+".docx")
	docx := doc.NewDoc()
	// 获取文件，并输出它们的名字
	err := filepath.Walk(tmpDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				errMsg <- err.Error()
				return err
			}
			if !info.IsDir() && file.IsImage(path) {
				docx.AddImageToWord(path)
				docx.SetImgWith(20)
			}
			return nil
		})
	if err != nil {
		errMsg <- err.Error()
		log.Println(err)
		return
	}
	ok := docx.Save(saveName)
	if !ok {
		errMsg <- "生成doc文档失败"
	}
}

func startCheckTmp() {
	remove()
	ticker := time.NewTicker(time.Second * 3 * 60 * 60)
	go func() {
		for { //循环
			<-ticker.C
			remove()
		}
	}()
}

func remove() {
	currentTime := time.Now()
	currentTime.Unix()
	n, _ := time.Parse("2006-01-02", currentTime.Format("2006-01-02"))
	tmpPath := filepath.Join(".", "tmp")
	files, err := ioutil.ReadDir(tmpPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		dirName := filepath.Join(tmpPath, f.Name())
		p, _ := time.Parse("2006-01-02", f.Name())
		if p.Unix() < n.Unix() {
			os.RemoveAll(dirName)
			log.Println("删除" + dirName)
		}
	}

}

