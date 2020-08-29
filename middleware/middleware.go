package middleware

import(
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context){
	log.Println("Start Middleware")
	authKey := c.GetHeader("Authorization")
	if authKey != "November 10, 2009"{
		c.JSON(http.StatusUnauthorized,"NOOOOOOOOOOOO")
		c.Abort()
		return
	}
	c.Next()
	log.Println("End Middleware")
}