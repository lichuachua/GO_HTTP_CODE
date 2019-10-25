package main

import (
	"net/http"
	"time"
)
func main()  {
	time1 := time.Now().Unix()
	url := "https://www.jianshu.com/p/f4863c11c13e"
	http.Get(url)
	println(url)
	for i:=0;i>=0 ;i++  {
		if (time.Now().Unix()-time1)>1 {
			http.Get(url)
			println(url)
			time1 = time.Now().Unix()
		}
	}
}