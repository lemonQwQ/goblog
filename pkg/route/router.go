package route

import (
	"github.com/gorilla/mux"
)

// Router 路由对象
var Router *mux.Router

// Initialize 初始化路由
func Initialize() {
	Router = mux.NewRouter()
}

// Name2URL 通过路由名称来获取URL
func Name2URL(routerName string, pairs ...string) string {
	url, err := Router.Get(routerName).URL(pairs...)
	if err != nil {
		// checkError(err)
		return ""
	}
	return url.String()
}
