package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"runtime"
	
	"bni.co.id/swift-filtering/config"
	"bni.co.id/swift-filtering/controllers"
	"bni.co.id/swift-filtering/database"
	"bni.co.id/swift-filtering/scheduler"
)

func main() {
	runtime.GOMAXPROCS(8)

	envars := getEnvVars()

	// "sqlserver://sa:Uat46@10.70.152.54:1433?database=SWIFT_FILTERING"
	database.SetupFactory(envars["swift.db.driver"], envars["swift.db.url"])
	defer database.Close()

	config.LoadSettings()
	scheduler.StartSendResponseJob()
	scheduler.StartClearCacheJob()

	router := gin.Default()
	filterCtrl := controllers.NewFilter()
	cacheCtrl := controllers.NewCache()

	v1 := router.Group("/")
	{
		// v1.Use(security.NewAuthMiddleware())

		v1.POST("/filter/check", filterCtrl.CheckFilter)
		v1.GET("/cache/clear-rules", cacheCtrl.ClearRules)
		v1.GET("/cache/clear-settings", cacheCtrl.ClearSettings)
	}

	router.Run()
}

func getEnvVars() config.Settings {
	var envars = config.Settings{}

	dbURL := os.Getenv("SWIFT_DB_URL")
	dbDriver := os.Getenv("SWIFT_DB_DRIVER")

	if dbURL == "" || dbDriver == "" {
		panic("Database configuration is not configured properly. Please configure DB connection first.")
		os.Exit(0)
	}

	envars["swift.db.url"] = dbURL
	envars["swift.db.driver"] = dbDriver

	return envars
}
