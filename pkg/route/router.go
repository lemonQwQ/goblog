package route

import (
	"goblog/pkg/config"
	"goblog/pkg/logger"
	"net/http"

	"github.com/gorilla/mux"
)

var route *mux.Router

func SetRoute(r *mux.Router) {
	route = r
}

// Name2URL 通过路由名称来获取URL
func Name2URL(routerName string, pairs ...string) string {
	url, err := route.Get(routerName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}
	return config.GetString("app.url") + url.String()
}

// GetRouteVariable 获取 URI 路由参数
func GetRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}
