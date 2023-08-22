package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func Authorize(c *gin.Context) {
	err := c.Request.ParseForm()
	if err != nil {
		c.String(http.StatusBadRequest, "bad request")
		c.Abort()
		return
	}

	userID := GetUserIdWhenAuth(c)
	if userID == 0 {
		c.Redirect(http.StatusPermanentRedirect, "/login")
		c.Abort()
		return
	}

	c.Redirect(http.StatusMovedPermanently, "../enterprises?manager="+strconv.Itoa(int(userID)))
}
