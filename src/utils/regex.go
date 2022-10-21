package utils

import "regexp"

// 正则表达式
const (
	// PHONE_REGEX 手机号
	PHONE_REGEX = "^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\\d{8}$"
	// EMAIL_REGEX 邮箱
	EMAIL_REGEX = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$"
	// PASSWORD_REGEX 密码
	PASSWORD_REGEX = "^\\w{4,32}$"
	// VERIFY_CODE_REGEX 验证码
	VERIFY_CODE_REGEX = "^[a-zA-Z\\d]{6}$"
)

func IsPhoneValid(phone string) bool {
	compiler := regexp.MustCompile(PHONE_REGEX)
	return compiler.MatchString(phone)
}
