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
	"sync"
	"sync/atomic"
	"time"

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

func (c *WriteCounter) Add(n uint64) {
	newCurrent := atomic.AddUint64(&c.Current, n)
	persentage := float64(newCurrent) / float64(c.Total) * 100
	runtime.EventsEmit(c.Ctx, "download_progress", persentage)
}

// Greet returns a greeting for the given name
func (a *App) DownloadDemo(url string) string {
	//发起下载请求
	resp, err := http.Head(url)
	if err != nil {
		return "获取文件信息失败：" + err.Error()
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
		return "创建文件失败"
	}
	defer out.Close()

	//获取文件大小
	filesize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if filesize <= 0 {
		return "获取文件大小失败" + err.Error()
	}

	//配置
	threadCount := 10
	partSize := filesize / threadCount

	//获取进度
	counter := &WriteCounter{
		Total: uint64(filesize),
		Ctx:   a.ctx,
	}

	var wg sync.WaitGroup

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		lastBytes := uint64(0)

		for {
			select {
			case <-ticker.C:
				currentBytes := atomic.LoadUint64(&counter.Current)
				speed := currentBytes - lastBytes

				lastBytes = currentBytes

				runtime.EventsEmit(a.ctx, "download_speed", speed)

				if currentBytes >= counter.Total {
					return
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()

	//开始下载
	for i := 0; i < threadCount; i++ {
		//计算起始位置和结束位置
		start := i * partSize
		end := start + partSize - 1

		if i == threadCount-1 {
			end = filesize - 1
		}

		wg.Add(1)

		go func(threadID int, startPos, endPos int) {
			defer wg.Done()
			a.downloadPart(url, out, startPos, endPos, counter)
		}(i, start, end)
	}

	wg.Wait()

	return fmt.Sprintf("下载完成! 大小：%d MB", filesize/1024/1024)
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

func (a *App) downloadPart(url string, file *os.File, start, end int, counter *WriteCounter) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("下载失败: %v\n", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 32*1024)
	currentOffset := int64(start)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			file.WriteAt(buf[:n], currentOffset)
			//更新偏移
			currentOffset += int64(n)
			//更新总进度
			counter.Add(uint64(n))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("读取错误： %v\n", err)
			break
		}
	}
}
