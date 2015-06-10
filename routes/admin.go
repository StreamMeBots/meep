package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/StreamMeBots/meep/pkg/user"
)

func getUsers(ctx *gin.Context) {
	usrs, err := user.Users()
	if err != nil {
		ctx.JSON(500, map[string]string{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"users": usrs,
	})
	return
}
