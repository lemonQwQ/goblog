package logger

import (
	"log"
)

// LogError 记录错误日志
func LogError(err error) {
	if err != nil {
		// log.Fatal(err)
		log.Printf("Server abort! Cause: %v \n", err)
	}
}
