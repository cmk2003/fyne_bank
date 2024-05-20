package main

//
//import (
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"sync"
//)
//
//func main() {
//	// 定义基础URL
//	baseURL := "http://eci-2ze3cmcfbu4hnt5lkudl.cloudeci1.ichunqiu.com/admin"
//
//	// 扩展字符集以包括小写字母和数字
//	charSet := "abcdefghijklmnopqrstuvwxyz1234567890"
//
//	// 使用sync.WaitGroup来等待所有goroutine完成
//	var wg sync.WaitGroup
//
//	// 遍历长度从1到4的所有字符串
//	for length := 1; length <= 4; length++ {
//		wg.Add(1)
//		go func(l int) {
//			defer wg.Done()
//			generateAndSend("", l, charSet, baseURL)
//		}(length)
//	}
//
//	wg.Wait()
//}
//
//func generateAndSend(s string, level int, charSet, baseURL string) {
//	if level == 0 {
//		fullURL := baseURL + s + ".php"
//		sendRequest(fullURL)
//		return
//	}
//	for _, char := range charSet {
//		generateAndSend(s+string(char), level-1, charSet, baseURL)
//	}
//}
//
//func sendRequest(url string) {
//	fmt.Println(url)
//	response, err := http.Get(url)
//	if err != nil {
//		fmt.Println("Error fetching URL: ", err)
//		return
//	}
//	defer response.Body.Close()
//
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		fmt.Println("Error reading response: ", err)
//		return
//	}
//	if response.Status == "404 Not Found" {
//		return
//	}
//	fmt.Printf("URL: %s\nStatus: %s\nBody: %s\n\n", url, response.Status, string(body))
//}
