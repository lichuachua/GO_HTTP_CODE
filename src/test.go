package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)



func test_HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024*4)
	for  {
		n,_ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		result += string(buf[:n])
	}
	return
}

func test_SpiderOneJoy(url string) (title, content string, err error) {

	result, err1 := test_HttpGet(url)
	if err1 != nil {
		err = err1
		return
	}

	//取标题
	re1 := regexp.MustCompile(`<h1 class="post-title">(?s:(.*?))</h1>`)
	if re1 == nil {
		err = fmt.Errorf("%s","regexp.MustCompile err")
		return
	}

	tmpTitle := re1.FindAllStringSubmatch(result, -1)
	for _, data := range tmpTitle{
		title = data[1]
		title = strings.Replace(title, "\r","",-1)
		title = strings.Replace(title, "\n","",-1)
		title = strings.Replace(title, "","",-1)
		title = strings.Replace(title, "\t","",-1)
		break
	}

	re2 := regexp.MustCompile(`<section class="post-content">(?s:(.*?))</section>`)
	if re2 == nil {
		err = fmt.Errorf("%s","regexp.MustCompile err")
		return
	}

	tmpContent := re2.FindAllStringSubmatch(result, -1)
	for _, data := range tmpContent{
		content = data[1]
		content = strings.Replace(content, "\r","",-1)
		content = strings.Replace(content, "\n","",-1)
		content = strings.Replace(content, "","",-1)
		content = strings.Replace(content, "\t","",-1)
		content = strings.Replace(content, "<p>","",-1)
		content = strings.Replace(content, "</p>","",-1)
		content = strings.Replace(content, "<br>","\n",-1)
		break
	}
	return
}

func test_StoreJoyToFile(i int, fileTitle , fileContent []string) {
	f, err := os.Create(strconv.Itoa(i)+".txt")
	if err != nil {
		fmt.Println("os.Create err = ",err)
		return
	}
	defer f.Close()

	n := len(fileTitle)
	for i := 0; i < n ; i++ {
		f.WriteString(fileTitle[i]+"\n")
		f.WriteString(fileContent[i]+"\n")
		f.WriteString("\n")
		f.WriteString("\n")
	}



}


func test_SpiderPage(i int, page chan int) {

	//http://duanziwang.com/page/1/
	url := "http://duanziwang.com/page/"+strconv.Itoa(i)
	fmt.Printf("正在爬取第%d个页面\n",i)

	result, err := test_HttpGet(url)
	if err != nil {
		fmt.Println("CSDN_HttpGet err = ", err)
		return
	}

	//取    <h1 class="dp-b"><a href="   段子url地址  "=====>>>正则表达式

	re := regexp.MustCompile(`<h1 class="post-title"><a href="(?s:(.*?))">`)
	if re == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}

	//取关键信息===>>每一页的段子链接组成一个切片
	joyUrl := re.FindAllStringSubmatch(result,-1)

	fileTitle := make([]string,0)
	fileContent := make([]string,0)

	//取网址
	for _,data := range joyUrl{
		title, content, err := test_SpiderOneJoy(data[1])
		if err != nil {
			fmt.Println("test_SpiderOneJoy err = ",err)
			continue
		}
		fileTitle = append(fileTitle, title)
		fileContent = append(fileContent, content)

	}
	test_StoreJoyToFile(i,fileTitle,fileContent)

	page <- i //写内容，写num
}

func test_Dowork(start, end int) {
	fmt.Printf("准备爬取第%d到%d页的网址\n",start,end)

	page := make(chan int)


	for i := start; i <= end ; i++ {
		go test_SpiderPage(i,page)
	}

	for i := start; i<= end; i++ {
		fmt.Printf("第%d个页面爬取完成\n",  <-page)
	}
}


func main()  {
	var start, end int
	fmt.Println("输入爬取的起始页(>=1)")
	fmt.Scan(&start)
	fmt.Println("输入爬取的起始页(>=start)")
	fmt.Scan(&end)

	test_Dowork(start,end)
}


