package controllers

import (
	"fmt"
	"net/http"
)

// PageController 处理静态页面
type PageController struct {
}

// Home 首页
func (*PageController) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

// About 关于我们页面
func (*PageController) About(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "博客用以记录笔记，如有反馈请联系："+"<a href=\"#\">tt</a>")
}

// NotFound 404 页面
func (*PageController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到:(</h1>")
}
