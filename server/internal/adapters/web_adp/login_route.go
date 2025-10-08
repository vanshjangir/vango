package web_adp

import "github.com/gin-gonic/gin"

type loginData struct {
	Type        string `json:"type"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Credentials string `json:"credential"`
}

func (gh *GinHandler) login(ctx *gin.Context) {
	var req loginData
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	switch req.Type {
	case "google":
		username, err := gh.us.LoginByGoogle(req.Credentials)
		if err != nil {
			ctx.JSON(403, gin.H{"error": "Login failed"})
		} else {
			ctx.JSON(200, gin.H{
				"message":  "Login successful",
				"username": username,
				"token":    req.Credentials,
			})
		}

	case "guest":
		username, token, err := gh.us.LoginAsGuest(req.Credentials)
		if err != nil {
			ctx.JSON(403, gin.H{"error": "Login failed"})
		} else {
			ctx.JSON(200, gin.H{
				"message":  "Login successful",
				"username": username,
				"token":    token,
			})
		}
	}
}
