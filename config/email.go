package config

import "goblog/pkg/config"

func init() {
	config.Add("email", config.StrMap{
		"from":    config.Env("EMAIL_FROM", ""),
		"sender":  config.Env("EMAIL_SENDER", ""),
		"subject": config.Env("EMAIL_SUBJECT", ""),
		"pwd":     config.Env("EMAIL_PWD", ""),
		"host":    config.Env("EMAIL_HOST", ""),
		"port":    config.Env("EMAIL_PORT", ""),
		"type":    "text/html",
	})
}
