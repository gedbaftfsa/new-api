package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"new-api/common"
	"new-api/middleware"
	"new-api/model"
	"new-api/router"
)

func main() {
	common.SetupLogger()
	common.SysLog("New API starting...")

	// Initialize database
	err := model.InitDB()
	if err != nil {
		common.FatalLog("Failed to initialize database: " + err.Error())
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.SysError("Failed to close database: " + err.Error())
		}
	}()

	// Initialize options from database
	err = model.InitOptionMap()
	if err != nil {
		common.FatalLog("Failed to initialize options: " + err.Error())
	}

	// Initialize memory cache
	common.InitTokenEncoders()

	// Set Gin mode — default to release unless explicitly set to debug
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "debug" {
		gin.SetMode(gin.DebugMode)
		common.SysLog("Running in DEBUG mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	middleware.SetupGlobalMiddleware(server)

	// Register routes
	router.SetRouter(server)

	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}

	common.SysLog(fmt.Sprintf("Server listening on port %s", port))
	err = server.Run(":" + port)
	if err != nil {
		common.FatalLog("Failed to start server: " + err.Error())
	}
}
