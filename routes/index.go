package routes

import (
	"alquran/authorization"
	"alquran/handler"
	"alquran/helpers"
	"alquran/kedaihelpers"
	"alquran/surah"
	"alquran/users"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Routing(router *gin.Engine, dbs kedaihelpers.DBStruct, initGorm *gorm.DB) {

	// repository
	userRepository := users.NewRepository(initGorm)
	surahRepository := surah.NewRepository(dbs)

	// services
	userServices := users.NewServices(userRepository)
	authServices := authorization.NewServices()
	surahServices := surah.NewServices(surahRepository)
	// handler
	userHandler := handler.NewUserHandler(userServices, authServices)
	surahHandler := handler.NewSurahHandler(surahServices)

	versioning := router.Group("api/v1")
	versioning.Any("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "OK",
			"Message": "Welcome to " + viper.GetString("appName"),
		})
	})
	authRouter := versioning.Group("auth")
	{
		authRouter.POST("/login/", userHandler.Login)
		authRouter.POST("/register/", userHandler.RegisterUser)
		authRouter.Use(AuthMiddleware(authServices, userServices))
		authRouter.POST("/refresh-token/", userHandler.RefreshToken)

	}
	surahRouter := versioning.Group("surah")
	{
		surahRouter.POST("/all/", surahHandler.GetSurah)
		surahRouter.POST("/detail/", surahHandler.GetDetailSurah)
	}
}

func AuthMiddleware(authServices authorization.Services, userServices users.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helpers.APIResponse("Unauthorized Access ", http.StatusUnauthorized, "error", "Yout Not Have Permission To Access This Site! Please Login or Exit")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		token, err := authServices.ValidateToken(tokenString)
		if err != nil {
			response := helpers.APIResponse("Unauthorized Access", http.StatusUnauthorized, "error", "Yout Not Have Permission To Access This Site! Please Login or Exit")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claims, ok := token.Claims.(*authorization.JWTClaim)
		if !ok || !token.Valid {
			response := helpers.APIResponse("Unauthorized Access", http.StatusUnauthorized, "error", "Yout Not Have Permission To Access This Site! Please Login or Exit")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if claims.ExpiresAt < time.Now().Local().Unix() {
			response := helpers.APIResponse("Token Expired", http.StatusUnauthorized, "error", "Yout Not Have Permission To Access This Site! Please Login or Exit")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("current", claims)
		c.Next()
	}
}
