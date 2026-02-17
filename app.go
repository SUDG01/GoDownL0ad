package main

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"

	"strconv"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) DownloadDemo(url string) string {
	//发起下载请求
	resp, err := http.Get(url)
	if err != nil {
		return "下载失败！"
	}
	defer resp.Body.Close()

	//寻找文件名
	contentDisposition := resp.Header.Get("Content-Disposition")
	finalFileName := ""

	if contentDisposition != "" {
		_, params, _ := mime.ParseMediaType(contentDisposition)
		finalFileName = params["filename"]
	}

	if finalFileName == "" {
		finalFileName = a.ParseFileName(url)
	}

	//创建文件
	filename := finalFileName
	out, err := os.Create(filename)
	if err != nil {
		return "创建文件失败，请检查文件权限."
	}
	defer out.Close()

	//获取文件大小
	filesize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	//获取进度
	counter := &WriteCounter{
		Total: uint64(filesize),
		Ctx:   a.ctx,
	}

	//开始下载
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		return "写入失败" + err.Error()
	}
	return "下载完成！文件位置：" + filename
}

type WriteCounter struct {
	Total   uint64          //文件总大小
	Current uint64          //当前已下载大小
	Ctx     context.Context //与前端通信
}

// 进度写入函数
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += uint64(n)

	//百分比计算
	percentage := float64(wc.Current) / float64(wc.Total) * 100

	//传递事件
	runtime.EventsEmit(wc.Ctx, "download_progress", percentage)
	return n, nil
}

func (a *App) ParseFileName(downloadUrl string) string {
	// 1. 解析 URL 结构 (处理 http://... 部分)
	u, err := url.Parse(downloadUrl)
	if err != nil {
		fmt.Println("URL 解析失败:", err)
		return "unknown_file"
	}

	// 2. 获取路径的最后一段
	// u.Path 会自动忽略 ?后面的参数 (比如 ?token=123)
	// path.Base("/media/gallery/SU%E5...png") -> "SU%E5...png"
	encodedName := path.Base(u.Path)

	// 3. URL 解码 (关键步骤！把 %xx 变回中文)
	decodedName, err := url.QueryUnescape(encodedName)
	if err != nil {
		fmt.Println("解码失败:", err)
		return encodedName // 如果解码失败，就返回带 % 的原名
	}

	return decodedName
}
