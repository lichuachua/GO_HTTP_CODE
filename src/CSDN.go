package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func CSDN_HttpGet(url string) (result string, err error) {

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

func CSDN_SpiderOneJoy(url string) (title string, err error) {
	fmt.Printf("%s\n",url)

	result, err1 := CSDN_HttpGet(url)
	if err1 != nil {
		err = err1
		return
	}

	//取标题
	re1 := regexp.MustCompile(`<h1 class="title-article">(?s:(.*?))</h1>`)
	if re1 == nil {
		err = fmt.Errorf("%s","regexp.MustCompile err")
		return
	}
	tmpTitle := re1.FindAllStringSubmatch(result, 1)

	for _, data := range tmpTitle{
		title = data[1]
		fmt.Println(title)
		break
	}

	return
}

func CSDN_StoreJoyToFile(i int, fileTitle []string) {
	f, err := os.Create(strconv.Itoa(i)+".txt")
	if err != nil {
		fmt.Println("os.Create err = ",err)
		return
	}
	defer f.Close()

	n := len(fileTitle)
	for i := 0; i < n ; i++ {
		f.WriteString(fileTitle[i]+"\n")
		f.WriteString("\n")
		f.WriteString("\n")
	}



}


func CSDN_SpiderPage(i int, page chan int) {

	//https://blog.csdn.net/qq_42410605/article/list/3?
	url := "https://blog.csdn.net/qq_42410605/article/list/"+strconv.Itoa(i)+"?"
	fmt.Printf("正在爬取第%d个页面%s\n",i,url)

	result, err := CSDN_HttpGet(url)
	if err != nil {
		fmt.Println("CSDN_HttpGet err = ", err)
		return
	}

	//取    <h1 class="dp-b"><a href="   博客url地址  "=====>>>正则表达式

	re := regexp.MustCompile(`<h4 class="">
        <a href="(?s:(.*?))" target="_blank">`)
	if re == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}
	//取关键信息===>>每一页的段子链接组成一个切片
	joyUrl := re.FindAllStringSubmatch(result,-1)

	fileTitle := make([]string,0)

	//取网址
	for _,data := range joyUrl{
		title, err := CSDN_SpiderOneJoy(data[1])
		if err != nil {
			fmt.Println("CSDN_SpiderOneJoy err = ",err)
			continue
		}
		fileTitle = append(fileTitle, title)

	}

	CSDN_StoreJoyToFile(i,fileTitle)

	page <- i //写内容，写num
}


func CSDN_Dowork(start, end int) {
	fmt.Printf("准备爬取第%d到%d页的网址\n",start,end)

	page := make(chan int)


	for i := start; i <= end ; i++ {
		go CSDN_SpiderPage(i,page)
	}

	for i := start; i<= end; i++ {
		fmt.Printf("第%d个页面爬取完成\n",  <-page)
	}
}


func main()  {
	time1 := time.Now().Unix()
	CSDN_Dowork(1,10)
	time2 := time.Now().Unix()
	for i:=0;i>=0 ;i++  {
		if (time.Now().Unix()-time1)>180 {
			CSDN_Dowork(1,10)
			time1 = time.Now().Unix()
		}
	}

	fmt.Printf("一共花费%v秒\n",time2-time1)
}
