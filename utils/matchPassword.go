package utils

import (
	"errors"
	"github.com/dlclark/regexp2"
)

func MatchStr(str string) error {
	expr := `^(?![0-9a-zA-Z]+$)(?![a-zA-Z!@#$%^&*]+$)(?![0-9!@#$%^&*]+$)[0-9A-Za-z!@#$%^&*]{8,16}$`
	reg, _ := regexp2.Compile(expr, 0)
	m, _ := reg.FindStringMatch(str)
	if m != nil {
		return nil
	}
	return errors.New("密码包含至少一位数字，字母和特殊字符,且长度8-16")
}
