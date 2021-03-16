package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/config"
	"goblog/pkg/flash"
	"goblog/pkg/logger"
	PWD "goblog/pkg/password"
	"goblog/pkg/session"
	"goblog/pkg/view"
	"net/http"

	"gopkg.in/gomail.v2"
)

// AuthController 处理静态页面
type AuthController struct {
}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	// 2. 表单规则
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		_user.Create()

		if _user.ID > 12 {
			// 登录用户并跳转到首页
			flash.Success("恭喜您注册成功！")
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "注册失败，请联系管理员")
		}
	}
}

// Login 显示登录表单
func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin 处理登录表单提交
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// 初始化表单数据
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// fmt.Fprint(w, email, password)
	// 尝试登录
	if err := auth.Attempt(email, password); err == nil {
		flash.Success("欢迎回来！")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Password": password,
		}, "auth.login")
	}
}

// Logout 退出登录
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	flash.Success("您已退出登录")

	http.Redirect(w, r, "/", http.StatusFound)
}

// Retrieve 显示找回密码页面
func (*AuthController) Retrieve(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.retrieve")
}

// DoRetrievc 处理找回密码表单提交
func (*AuthController) DoRetrieve(w http.ResponseWriter, r *http.Request) {
	to := r.PostFormValue("email")

	err := auth.VerifyEmail(to)
	if err == nil {
		from := config.GetString("email.from")
		// 配置邮件信息
		m := gomail.NewMessage()
		m.SetAddressHeader("From", from, config.GetString("email.sender"))
		m.SetHeader("To", to)
		m.SetHeader("Subject", config.GetString("email.subject"))
		// m.Embed()
		m.SetBody(config.GetString("email.type"), "123456")

		// 发生邮件
		d := gomail.NewDialer(config.GetString("email.host"), config.GetInt("email.port"), from, config.GetString("email.pwd"))

		if err := d.DialAndSend(m); err != nil {
			logger.LogError(err)
			flash.Danger(config.GetString("email.host") + config.GetString("email.port") + from + config.GetString("email.pwd"))
			http.Redirect(w, r, "/", http.StatusFound)
		}
		view.RenderSimple(w, view.D{
			"Email": to,
		}, "auth.verification")
	} else {
		session.Flush()
		view.RenderSimple(w, view.D{
			"Error": err.Error(),
			"Email": to,
		}, "auth.retrieve")
	}
}

// ModifyPwd 显示修改密码页面
func (*AuthController) ModifyPwd(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.modifypwd")
}

// DoModifyPwd 处理修改密码表单提交
func (*AuthController) DoModifyPwd(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	passwordConfirm := r.PostFormValue("password_confirm")

	// 2. 表单规则
	errs := requests.ValidatePwd(password, passwordConfirm)

	if len(errs) > 0 {
		view.RenderSimple(w, view.D{
			"Errors":          errs,
			"Email":           email,
			"Password":        password,
			"PasswordConfirm": passwordConfirm,
		}, "auth.modifypwd")
	} else {
		_user, _ := user.GetByEmail(email)
		err := _user.Update(PWD.Hash(password))
		if err == nil {
			// 登录用户并跳转到首页
			flash.Success("恭喜您修改成功！")
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "修改失败，请联系管理员")
		}
	}
}

// Verification 显示邮箱验证页面
func (*AuthController) Verification(w http.ResponseWriter, r *http.Request) {

}

// DoVerification 处理邮箱验证表单提交
func (*AuthController) DoVerification(w http.ResponseWriter, r *http.Request) {

}
