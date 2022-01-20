package handle

import (
	"edgelog/routers/common"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleUser(group *gin.RouterGroup) {
	group.GET("/debug", debug)
	group.POST("/login", login)
	group.GET("/logout", logout)
	group.GET("/username", username)
}

// @Summary  login
// @Tags  user
// @Param username formData string true " username "
// @param password formData string true " password "
// @Router /api/v1/user/login [post]
func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != common.USERNAME {
		c.JSON(500, gin.H{
			"message": "username error",
		})
		return
	}
	if password != common.PASSWORD {
		c.JSON(500, gin.H{
			"message": "password error",
		})
		return
	}
	tokenStr, err := common.GetOrAddToken(common.GetTokenKey(c), common.Token{
		Expire: time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"message": tokenStr,
		})
	}
}

// @Summary  logout
// @Tags  user
// @Router /api/v1/user/logout [get]
func logout(c *gin.Context) {
	common.TokenMap.Delete(common.GetTokenKey(c))
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func username(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": common.USERNAME,
	})
}

func debug(c *gin.Context) {
	result := make(map[string]interface{})
	common.TokenMap.Range(func(key, value interface{}) bool {
		result[key.(string)] = value
		return true
	})
	c.JSON(200, result)
}
