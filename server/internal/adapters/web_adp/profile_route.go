package web_adp

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vanshjangir/vango/server/internal/domain"
)

func (gh *GinHandler) profile(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		ctx.JSON(400, gin.H{"error": "username not provided"})
		return
	}
	userProfile, err := gh.us.GetUserProfile(username)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "user not found"})
		return
	}

	reviews, err := gh.gs.RecentGames(username, 10)
	if err != nil {
		log.Println("profile:", err)
	} else {
		userProfile.RecentGames = make([]domain.RecentGames, len(reviews))
		for i, g := range reviews {
			userProfile.RecentGames[i].Gameid = g.Gameid
			userProfile.RecentGames[i].CreatedAt = g.CreatedAt
			if g.WhiteName == username {
				userProfile.RecentGames[i].Opponent = g.BlackName
				if g.Winner == domain.WhiteColor {
					userProfile.RecentGames[i].Result = "Won"
					userProfile.Wins += 1
				} else {
					userProfile.RecentGames[i].Result = "Lost"
					userProfile.Losses += 1
				}
			} else {
				userProfile.RecentGames[i].Opponent = g.WhiteName
				if g.Winner == domain.BlackColor {
					userProfile.RecentGames[i].Result = "Won"
					userProfile.Wins += 1
				} else {
					userProfile.RecentGames[i].Result = "Lost"
					userProfile.Losses += 1
				}
			}
		}
	}
	userProfile.GamesPlayed = userProfile.Wins + userProfile.Losses

	ctx.JSON(200, userProfile)
}
