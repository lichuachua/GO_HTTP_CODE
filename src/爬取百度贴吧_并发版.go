package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

//并发爬取网页内容
func ConHttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}

	defer resp.Body.Close()

	//读取网页body内容
	buf := make([]byte, 1024*4)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			//读取结束，或者出问题
			fmt.Println("resp.Body.Read err = ", err)
			break
		}
		result += string(buf[:n])
	}
	return
}

func SpiderPage(i int,  page chan<- int){
	url := "https://tieba.baidu.com/f?kw=王者荣耀&ie=utf-8&cid=&tab=corearea&pn=" +
	//将页面值转换为字符串类型（(i-1)*50是int类型）
		strconv.Itoa((i-1)*50)
	fmt.Println("正在爬取第%d页", i)

	//2)爬  （将所有的网站内容全部爬下来）
	result, err := ConHttpGet(url)
	if err != nil {
		fmt.Println("HttpGet err = ", err)
		return
	}

	//把爬取的内容写入一个文件
	fileName := strconv.Itoa(i) + ".html"
	f, err1 := os.Create(fileName)
	if err1 != nil {
		fmt.Println("os.Create err1 = ", err1)
		return
	}
	f.WriteString(result) //写内容
	f.Close()             //关闭文件
	page <- i
}

func ConDoWork(start, end int) {
	fmt.Printf("正在爬取%d到%d的页面", start, end)
	fmt.Println()

	//使用管道
	page := make(chan int)
	//明确目标（要知道你准备在哪个范围或者网站上搜索）
	//https://tieba.baidu.com/f?kw=王者荣耀&ie=utf-8&cid=&tab=corearea&pn=0    //下一页要+50
	//迭代爬取内容
	for i := start; i <= end; i++ {
		//go 一个进程
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d个页面爬取完成\n", <-page)
	}
}

func main()  {
	var start, end int
	fmt.Println("输入爬取的起始页(>=1)")
	fmt.Scan(&start)
	fmt.Println("输入爬取的终止页(>=start)")
	fmt.Scan(&end)

	ConDoWork(start,end)
}

