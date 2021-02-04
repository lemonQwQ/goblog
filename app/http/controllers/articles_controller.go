package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
	BaseController
}

// Show 文章详情页面
func (ac *ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1.获取URL参数
	id := route.GetRouteVariable("id", r)
	// 2.读取对应的文章数据
	article, err := article.Get(id)
	// 3.如果出现错误
	if err != nil {
		ac.ResponceForSQLError(w, err)
	} else {
		// ---  4. 读取成功，显示文章 ---
		view.Render(w, view.D{
			"Article":          article,
			"CanModifyArticle": policies.CanModifyArticle(article),
		}, "articles.show", "articles._article_meta")
	}
}

// Index 文章列表页
func (ac *ArticlesController) Index(w http.ResponseWriter, r *http.Request) {

	// fmt.Fprint(w, config.Get("app.name"))

	articles, err := article.GetAll()
	// fmt.Println("文章数据：", articles)
	if err != nil {
		ac.ResponceForSQLError(w, err)
	} else {
		// ---  2. 加载模板 ---
		view.Render(w, view.D{
			"Articles": articles,
		}, "articles.index", "articles._article_meta")
	}
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	currentUser := auth.User()
	_article := article.Article{
		Title:  r.PostFormValue("title"),
		Body:   r.PostFormValue("body"),
		UserID: currentUser.ID,
	}
	errors := requests.ValidateArticleForm(_article)

	// 检查是否有错误
	if len(errors) == 0 {
		_article.Create()
		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusFound)
			// fmt.Fprint(w, "插入成功，ID 为"+_article.GetStringID())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.create", "articles._form_field")
	}
}

// Edit 文章更新页面
func (ac *ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {
	// 1.获取URL参数
	id := route.GetRouteVariable("id", r)

	// 2.读取对应的文章数据
	_article, err := article.Get(id)

	if err != nil {
		ac.ResponceForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			view.Render(w, view.D{
				"Article": _article,
				"Errors":  view.D{},
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Update 更新文章
func (ac *ArticlesController) Update(w http.ResponseWriter, r *http.Request) {
	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)
	if err != nil {
		ac.ResponceForSQLError(w, err)
	} else {
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {

			_article.Title = r.PostFormValue("title")
			_article.Body = r.PostFormValue("body")

			errs := requests.ValidateArticleForm(_article)

			if len(errs) == 0 {
				rowsAffected, err := _article.Update()

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "500 服务器内部错误")
					return
				}

				if rowsAffected > 0 {
					showURL := route.Name2URL("articles.show", "id", id)
					http.Redirect(w, r, showURL, http.StatusFound)
				} else {
					fmt.Fprint(w, "您没有做任何更改！")
				}
			} else {
				view.Render(w, view.D{
					"Article": _article,
					"Errors":  errs,
				}, "articles.edit", "articles._form_field")
			}
		}
	}
}

// Delete 删除文章
func (ac *ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {

	id := route.GetRouteVariable("id", r)
	_article, err := article.Get(id)

	if err != nil {
		ac.ResponceForSQLError(w, err)
	} else {
		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResponseForUnauthorized(w, r)
		} else {
			rowsAffected, err := _article.Delete()

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			} else {
				if rowsAffected > 0 {
					indexURL := route.Name2URL("articles.index")
					http.Redirect(w, r, indexURL, http.StatusFound)
				} else {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "404 文章未找到")
				}
			}
		}
	}
}
