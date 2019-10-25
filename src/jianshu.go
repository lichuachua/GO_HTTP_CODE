package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

func JianShu_HttpGet1(url string) (result string, err error) {

	resp, err1 := http.Get("https://www.jianshu.com"+url)
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

func JianShu_HttpGet(url string) (result string, err error) {

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

func JianShu_SpiderOneJoy(url string) (title string, err error) {
	fmt.Printf("%s\n",url)

	result, err1 := JianShu_HttpGet1(url)
	if err1 != nil {
		err = err1
		return
	}

	//取标题 <h1 class="_1RuRku">(?s:(.*?))</h1>
	re1 := regexp.MustCompile(`<h1 class="_1RuRku">(?s:(.*?))</h1>`)
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

func JianShu_StoreJoyToFile(i int, fileTitle []string) {
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


func JianShu_SpiderPage(i int, page chan int) {


	url := "https://www.jianshu.com/u/aee4339a566f"
	fmt.Printf("正在爬取第%d个页面 %s\n",i,url)

	result, err := JianShu_HttpGet(url)
	if err != nil {
		fmt.Println("JianShu_HttpGet err = ", err)
		return
	}

	//取    <a class="title" target="_blank" href="(?s:(.*?))">=====>>>正则表达式

	re := regexp.MustCompile(`<a class="title" target="_blank" href="(?s:(.*?))">`)
	if re == nil {
		fmt.Println("regexp.MustCompile err")
		return
	}
	//取关键信息===>>每一页的段子链接组成一个切片
	joyUrl := re.FindAllStringSubmatch(result,-1)

	fileTitle := make([]string,0)

	//取网址
	for _,data := range joyUrl{
		title, err := JianShu_SpiderOneJoy(data[1])
		if err != nil {
			fmt.Println("JianShu_SpiderOneJoy err = ",err)
			continue
		}
		fileTitle = append(fileTitle, title)

	}

	JianShu_StoreJoyToFile(i,fileTitle)

	page <- i //写内容，写num
}

func JianShu_Dowork(start, end int) {
	fmt.Printf("准备爬取第%d到%d页的网址\n",start,end)

	page := make(chan int)


	for i := start; i <= end ; i++ {
		go JianShu_SpiderPage(i,page)
	}

	for i := start; i<= end; i++ {
		fmt.Printf("第%d个页面爬取完成\n",  <-page)
	}
}


func main()  {
	time1 := time.Now().Unix()
	JianShu_Dowork(1,1)
	time2 := time.Now().Unix()
	fmt.Printf("一共花费%v秒\n",time2-time1)
}
