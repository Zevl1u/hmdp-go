package utils

import "time"

const (
	USER_NICK_NAME_PREFIX = "user_"
	LOGIN_CODE_PREFIX     = "login:code:"
	LOGIN_CODE_TTL        = 5 * time.Minute
	LOGIN_USERDTO_TTL     = 30 * time.Minute

	AUTHORIZATION = "Authorization"

	CACHE_SHOP_PREFIX   = "cache:shop:"
	CACHE_SHOP_INFO_TTL = 30 * time.Minute
)
