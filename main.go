package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

var router = mux.NewRouter()

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "博客用以记录笔记，如有反馈请联系："+"<a href=\"#\">tt</a>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID: "+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errs := make(map[string]string)

	// 验证标题
	if title == "" {
		errs["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errs["title"] = "标题长度需介于 3-40"
	}

	// 验证内容
	if body == "" {
		errs["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errs["body"] = "内容长度需大于或等于 10 个字节"
	}

	if len(errs) == 0 {
		fmt.Fprintf(w, "验证通过！<br>")
		fmt.Fprintf(w, "title 的值为：%v <br>", title)
		fmt.Fprintf(w, "title 的长度为：%v <br>", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body 的值为：%v <br>", body)
		fmt.Fprintf(w, "body 的长度为：%v <br>", utf8.RuneCountInString(body))
	} else {
		/*html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<title>创建文章 —— 我的技术博客</title>
			<style type="text/css">.error{color: red;}</style>
		</head>
		<body>
			<form action="{{ .URL }}" method="post">
				<p><input type="text" name="title" value={{ .Title }}></p>
				{{ with .Errors.title }}
				<p class="error"> {{ . }} </p>
				{{ end }}
				<p><textarea name="body" cols="30" rows="10">{{ .Body }}</textarea></p>
				{{ with .Errors.body }}
				<p class="error"> {{ . }} </p>
				{{ end }}
				<p><button type="submit">提交</button></p>
			</form>
		</body>
		</html>
		`*/
		storeURL, _ := router.Get("articles.store").URL()

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errs,
		}

		// tmpl, err := template.New("create-form").Parse(html)
		tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")

		if err != nil {
			panic(err)
		}

		tmpl.Execute(w, data)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到:(</h1>")
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置标头
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		// 进行处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	/*html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>创建文章 —— 我的技术博客</title>
	</head>
	<body>
		<form action="%s" method="post">
			<p><input type="text" name="title"></p>
			<p><textarea name="body" cols="30" rows="10"></textarea></p>
			<p><button type="submit">提交</button></p>
		</form>
	</body>
	</html>
	`*/
	storeURL, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, data)
}

func main() {
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	// 自定义404页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取URL示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL: ", homeURL)
	articlesURL, _ := router.Get("articles.show").URL("id", "233")
	fmt.Println("articlesURL: ", articlesURL)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
