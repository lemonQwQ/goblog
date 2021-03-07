package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

// UserController 用户控制器
type UserController struct {
	BaseController
}

// Show 用户个人页面
func (uc *UserController) Show(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_user, err := user.Get(id)

	articles, pagerData, err := article.GetByUserIDCS(_user.GetStringID(), r, 3)
	
	if err != nil {
		uc.ResponceForSQLError(w, err)
	} else {
		
		if err != nil {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			view.Render(w, view.D{
				"Articles":  articles,
				"PagerData": pagerData,
			}, "articles.index", "articles._article_meta")
		}
	}
}
