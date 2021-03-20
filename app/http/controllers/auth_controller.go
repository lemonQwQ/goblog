package controllers

import (
	"fmt"
	"goblog/app/models/user"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/email"
	"goblog/pkg/flash"
	"goblog/pkg/password"
	PWD "goblog/pkg/password"
	"goblog/pkg/session"
	"goblog/pkg/util"
	"goblog/pkg/view"
	"net/http"
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
		code := util.GetRandom(6)
		err := email.Send(to, code)
		if err != nil {
			flash.Danger("内部错误，请稍后再试！")
			http.Redirect(w, r, "/", http.StatusFound)
		}
		session.Put("hash", password.Hash(code))
		session.Put("email", to)
		view.RenderSimple(w, view.D{}, "auth.verification")
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
	code := r.PostFormValue("vcode")
	hash := session.Get("hash")
	if password.CheckHash(code, hash.(string)) {
		view.RenderSimple(w, view.D{}, "auth.modifypwd")
	} else {
		flash.Danger("验证失败,请重新验证!")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// DoModifyPwd 处理修改密码表单提交
func (*AuthController) DoModifyPwd(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	password := r.PostFormValue("password")
	passwordConfirm := r.PostFormValue("password_confirm")

	// 2. 表单规则
	errs := requests.ValidatePwd(password, passwordConfirm)

	if len(errs) > 0 {
		view.RenderSimple(w, view.D{
			"Errors":          errs,
			"Password":        password,
			"PasswordConfirm": passwordConfirm,
		}, "auth.modifypwd")
	} else {
		email := session.Get("email")
		session.Forget("email")
		session.Forget("hash")
		_user, _ := user.GetByEmail(email.(string))
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

// Verification 处理邮箱验证表单提交
func (*AuthController) Verification(w http.ResponseWriter, r *http.Request) {

}

// DoVerification 处理邮箱验证表单提交
func (*AuthController) DoVerification(w http.ResponseWriter, r *http.Request) {

}
