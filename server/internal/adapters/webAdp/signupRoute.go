package webAdp

import "github.com/gin-gonic/gin"

type signupData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (gh *GinHandler) signup(ctx *gin.Context) {
	var req signupData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
	}

	err := gh.us.Signup(req.Username, req.Email, req.Password)
	if err != nil {
		ctx.JSON(403, gin.H{"error": err})
	} else {
		ctx.JSON(200, gin.H{"message": "Signup successful"})
	}
}
