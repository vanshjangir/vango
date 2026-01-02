package web_adp

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/vango/server/internal/domain"
)

func (gh *GinHandler) findGame(ctx *gin.Context) {
	var wp domain.WaitingPlayer
	usernameItf, exists := ctx.Get("username")
	usertypeItf, exists := ctx.Get("usertype")
	
	username := usernameItf.(string)
	usertype := usertypeItf.(string)
	
	if !exists {
		log.Println("username and usertype not found")
		ctx.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	
	wp.Username = username
	wp.AddedAt = time.Now()
	if usertype == "guest" {
		wp.Rating = 400
	} else {
		user, err := gh.us.GetUser(username)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "user does not exists"})
			return
		}
		wp.Rating = user.Rating
	}
	
	
	if wsurl, err := gh.ms.Match(wp); err != nil {
		log.Println("Error in Match:", err)
		ctx.JSON(500, gin.H{"error": "Matching failed"})
	} else {
		ctx.JSON(200, gin.H{"wsurl": wsurl})
	}
}
