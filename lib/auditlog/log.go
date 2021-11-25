//
// @Author: QLeelulu
// @Date: 2021-11-25 13:00:44
// @LastEditors: QLeelulu
// @LastEditTime: 2021-11-25 13:00:44
// @FilePath: /aproxy/lib/auditlog/log.go
// @Description:

package auditlog

import (
	"fmt"
	"log"
	"os"

	"aproxy/module/auth"
)

var logger *log.Logger

// Init 初始化
func Init(logPath string) error {
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("file open error : %v", err)
	}
	logger = log.New(f, "", log.Ldate|log.Ltime)
	return nil
}

// AccessLog 打印访问审计日志
func AccessLog(u *auth.User, resources string) error {
	logger.Printf("%s[%s] %s", u.Name, u.Email, resources)
	return nil
}
