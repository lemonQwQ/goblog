package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
	// w.WriteHeader(http.StatusInternalServerError)
	// w.Header().Set("name", "my name is lemon")
	// fmt.Fprint(w, "请求路径为:"+r.URL.Path)
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //标头设置
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
	} else if r.URL.Path == "/about" {
		fmt.Fprint(w, "博客用以记录笔记，如有反馈请联系："+"<a href=\"https://www.baidu.com\">百度</a>")
	} else {
		fmt.Fprintf(w, "<h1>请求页面为找到:(</h1>")
	}
}
func main() {
	// ‘/’反斜杠代表任意路径
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":3000", nil)
}
