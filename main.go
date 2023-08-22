package main

import (
	"alquran/core"
	"alquran/routes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// ---- READ CONFIG JSON ---- //
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("app.conf")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	// ---- END READ CONFIG JSON ---- //

	// ---- Config Server Runner ---- //

	initGorm, err := core.InitGorm()
	router := gin.New()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dbs := core.DBConnect()
	defer dbs.Dbx.Close()

	// ---- Setup Config Cors ---- //
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"X-CSRF-Token", "Accept-Encoding", "Content-Length", "Origin", "Access-Control-Allow-Origin", "Authorization", "Content-Type", "x-requested-with", "referer"},
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin"},
		MaxAge:        12 * time.Hour,
	}))
	// ---- End Setup Config Cors ---- //

	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Next()
	})

	// ---- Setup Config S3 Cloud Storage ---- //
	S3 := core.ConnectAws()
	router.Use(func(c *gin.Context) {
		c.Set("S3", S3)
		c.Next()
	})
	// ---- End Setup Config S3 Cloud Storage ---- //

	// CronJOB
	// cronJob := cron.New()

	// cronJob.AddFunc("0 */12 * * *", func() {
	// fetch.RunnerFetching(dbs)
	// })
	// cronJob.Start()
	// ---- List Routing ---- //
	routes.Routing(router, dbs, initGorm)
	// ---- End routing ---- //

	tmphttpreadheadertimeout, _ := time.ParseDuration(viper.GetString("server.readheadertimeout") + "s")
	tmphttpreadtimeout, _ := time.ParseDuration(viper.GetString("server.readtimeout") + "s")
	tmphttpwritetimeout, _ := time.ParseDuration(viper.GetString("server.writetimeout") + "s")
	tmphttpidletimeout, _ := time.ParseDuration(viper.GetString("server.idletimeout") + "s")
	s := &http.Server{
		Addr:              ":" + viper.GetString("server.port"),
		Handler:           router,
		ReadHeaderTimeout: tmphttpreadheadertimeout,
		ReadTimeout:       tmphttpreadtimeout,
		WriteTimeout:      tmphttpwritetimeout,
		IdleTimeout:       tmphttpidletimeout,
	}

	fmt.Println("🚀 Server running on port:", viper.GetString("server.port"))
	s.ListenAndServe()
	// ---- End Config Server Runner ---- //

}
