package auth

import (
	"errors"
	"goblog/app/models/user"
	"goblog/pkg/session"

	"gorm.io/gorm"
)

func _getUID() string {
	_uid := session.Get("uid")
	uid, ok := _uid.(string)
	if ok && len(uid) > 0 {
		return uid
	}
	return ""
}

// User 获取登录用户信息
func User() user.User {
	uid := _getUID()
	if len(uid) > 0 {
		_user, err := user.Get(uid)
		if err == nil {
			return _user
		}
	}
	return user.User{}
}

// Attempt 尝试登录
func Attempt(email string, password string) error {
	// 根据 Email 获取用户
	_user, err := user.GetByEmail(email)

	// 出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("邮箱不存在或者密码错误")
		} else {
			return errors.New("内部错误，请稍后尝试")
		}
	}

	// 匹配密码
	if !_user.ComparePassword(password) {
		return errors.New("邮箱不存在或者密码错误")
	}

	// 登录用户，报错会话
	session.Put("uid", _user.GetStringID())

	return nil
}

// Verification 验证邮箱
func Verification(email string) error {
	// 根据 Email 验证用户是否存在
	if _, err := user.GetByEmail(email); err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("邮箱不存在")
		} else {
			return errors.New("内部错误，请稍后尝试")
		}
	}

	return nil
}

// Login 登录指定用户
func Login(_user user.User) {
	session.Put("uid", _user.GetStringID())
}

// Logout 退出用户
func Logout() {
	session.Forget("uid")
}

// Check 检测是否登录
func Check() bool {
	return len(_getUID()) > 0
}
