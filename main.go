package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
	// w.WriteHeader(http.StatusInternalServerError)
	// w.Header().Set("name", "my name is lemon")
	// fmt.Fprint(w, "请求路径为:"+r.URL.Path)
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //标头设置
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "<h1>请求页面为找到:(</h1>")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "博客用以记录笔记，如有反馈请联系："+"<a href=\"#\">tt</a>")
}

func main() {
	router := http.NewServeMux()
	// ‘/’反斜杠代表任意路径
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", aboutHandler)
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 4)[2]
		fmt.Fprint(w, "文章id = "+id)
	})
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w, "GET method request")
		case "POST":
			fmt.Fprint(w, "POST method request")
		}
	})
	http.ListenAndServe(":3000", router)
}
