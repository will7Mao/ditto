package main

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type DittoForm struct {
	Method      string
	Url         string
	Body        string
	Loop        int
	Concurrency int
	Headers     []HEADER
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

func DoDitto(request Request, headers []HEADER) {
	log.Println("Ditto开始,并发数：%v循环:%v", request.Concurrency, request.Loop)
	log.Println(request)
	log.Println(headers)
	// 创建等待组，以便在所有协程完成后退出程序
	var wg sync.WaitGroup
	wg.Add(request.Concurrency)

	// 创建HTTP请求客户端
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
		},
		Timeout: 10 * time.Second,
	}

	// 创建指定数量的并发协程
	for i := 0; i < request.Concurrency; i++ {
		go func() {
			defer wg.Done()

			for loop := 0; loop < request.Loop; loop++ {
				// 复用HTTP请求
				req, err := http.NewRequest(request.Method, request.Url, strings.NewReader(request.RequestBody))

				// 设置请求头
				for _, header := range headers {
					req.Header.Set(header.Key, header.Value)
				}
				if err != nil {
					log.Printf("Failed to create request: %v\n", err)
					return
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
				// 确保响应正常关闭
				defer resp.Body.Close()

				// 可以根据需要在这里处理响应，例如检查状态码、读取响应体等
			}
		}()
	}
	// 等待所有协程完成
	wg.Wait()
}
