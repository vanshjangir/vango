package web_adp

import (
	"log"

	"github.com/gin-gonic/gin"
)

type ChangeUsernameType struct {
	Username    string `json:"username"`
	Newusername string `json:"newusername"`
}

func (gh *GinHandler) changeUsername(ctx *gin.Context) {
	var cu ChangeUsernameType
	if err := ctx.ShouldBindJSON(&cu); err != nil {
		log.Println("JSON bind error: ", err)
		ctx.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	err := gh.us.ChangeUsername(cu.Username, cu.Newusername)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Error updating username"})
	} else {
		ctx.JSON(200, "Username changed successfully")
	}
}
