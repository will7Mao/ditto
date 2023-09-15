package main

import (
	"embed"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed ditto_front_end/dist/*
var static embed.FS

type DittoForm struct {
	Method      string   `json:"method" binding:"required"`
	Url         string   `json:"url" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	Loop        int      `json:"loop" binding:"required"`
	Concurrency int      `json:"concurrency" binding:"required"`
	Headers     []HEADER `json:"headers" binding:"required"`
}

func main() {
	// 在运行时将关闭调试信息和堆栈跟踪
	router := gin.Default()

	// 静态文件
	router.StaticFS("/static", http.FS(static))

	router.POST(
		"/api/v1/ditto", func(ctx *gin.Context) {
			log.Println("处理压测请求")
			var requestJson DittoForm
			if err := ctx.ShouldBindJSON(&requestJson); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Printf("requestJson: %v\n", requestJson)
			Ditto(
				Request{
					Url:         requestJson.Url,
					Method:      requestJson.Method,
					RequestBody: requestJson.Body,
					Concurrency: requestJson.Concurrency,
					Loop:        requestJson.Loop,
				},
				requestJson.Headers,
			)
			ctx.JSON(http.StatusOK, gin.H{"status": "ditto"})
		})

	// 当 API 不存在时，返回静态文件
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path                                   // 获取请求路径
		s := strings.Split(path, ".")                                // 分割路径，获取文件后缀
		prefix := "ditto_front_end/dist"                             // 前缀路径
		if data, err := static.ReadFile(prefix + path); err != nil { // 读取文件内容
			// 如果文件不存在，返回首页 index.html
			if data, err = static.ReadFile(prefix + "/index.html"); err != nil {
				c.JSON(404, gin.H{
					"err": err,
				})
			} else {
				c.Data(200, mime.TypeByExtension(".html"), data)
			}
		} else {
			// 如果文件存在，根据请求的文件后缀，设置正确的mime type，并返回文件内容
			c.Data(200, mime.TypeByExtension(fmt.Sprintf(".%s", s[len(s)-1])), data)
		}
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run()
}

type HEADER struct {
	Key   string
	Value string
}

type Request struct {
	Url         string
	Method      string
	RequestBody string
	Concurrency int
	Loop        int
}

/**
 * ditto
 */
func Ditto(request Request, headers []HEADER) {
	log.Println("Ditto开始,并发数：%v循环%v", request.Concurrency, request.Loop)
	log.Println(request)
	log.Println(headers)
	// 创建等待组，以便在所有协程完成后退出程序
	var wg sync.WaitGroup
	wg.Add(request.Concurrency)

	// 创建指定数量的并发协程
	for i := 0; i < request.Concurrency; i++ {
		go func() {
			defer wg.Done()

			for loop := 0; loop < request.Loop; loop++ {
				// 创建HTTP请求客户端
				client := &http.Client{
					Timeout: 10 * time.Second,
				}

				// 创建HTTP请求
				req, err := http.NewRequest(request.Method, request.Url, strings.NewReader(request.RequestBody))

				if err != nil {
					log.Printf("Failed to create request: %v\n", err)
					return
				}

				// 设置请求头
				for _, header := range headers {
					req.Header.Set(header.Key, header.Value)
				}

				// 发送HTTP请求
				resp, err := client.Do(req)
				if err != nil {
					log.Printf("Failed to send request: %v\n", err)
					return
				}

				// body, err := ioutil.ReadAll(resp.Body)
				// if err != nil {
				// 	fmt.Printf("Failed to read response body: %v\n", err)
				// 	return
				// }

				// fmt.Println("Response Body:", string(body))

				// 确保响应正常关闭
				resp.Body.Close()

				// 可以根据需要在这里处理响应，例如检查状态码、读取响应体等
			}
		}()
	}

	// 不等待所有协程完成
	// wg.Wait()
}
