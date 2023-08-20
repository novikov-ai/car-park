package auth

import "github.com/gin-gonic/gin"

const (
	secretKey = "secret"
	userIdKey = "userId"
)

var userInfo = map[string]map[string]interface{}{
	"ismirnov@example.com": {
		"secret": "456",
		"userId": int64(1),
	},
	"mgreen@example.com": {
		"secret": "123",
		"userId": int64(2),
	},
}

func GetUserIdWhenAuth(c *gin.Context) int64 {
	login := c.Request.FormValue("login")

	user, ok := userInfo[login]
	if !ok {
		return 0
	}

	secret, ok := user[secretKey]
	if !ok {
		return 0
	}

	password := c.Request.FormValue("password")

	if secret != password {
		return 0
	}

	userID := user[userIdKey]
	return userID.(int64)
}
