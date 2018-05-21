package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/richardsang2008/BotManager/utility"
	"github.com/spf13/viper"
	"time"

	"context"
	"github.com/richardsang2008/BotManager/controller"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/services"
	"io"
	"net/http"
	"os"
	"os/signal"
)

func GetUsers(c *gin.Context) {
	//	log.Debug("I am here")
	c.JSON(200, "hello world")
}

func setupRouter(router *gin.Engine)  {
	router.POST("/account/add", services.AddAccount)
	router.GET("/account/add", services.AddAccount)
	router.POST("/account/update", services.UpdateAccountBySpecificFields)
	router.GET("/account/request", services.GetAccountBySystemIdAndLevelAndMark)
	router.POST("/account/release", services.ReleaseAccount)
	//end of meet the old one
	router.POST("/ptcaccounts/accounts/v1", services.AddAccount)
	router.GET("/ptcaccounts/accounts/v1/id/:id", services.GetAccountById)
	router.GET("/ptcaccounts/accounts/v1/", services.GetAccountByUserName)
	//router.POST("/ptcaccounts/accounts/v1/lvl/:level", services.AddAccountWithLevelHandler(maxlevel))
	router.PATCH("/ptcaccounts/accounts/v1/release", services.ReleaseAccount)
	router.GET("/ptcaccounts/accounts/v1/request", services.GetAccountBySystemIdAndLevelAndMark)
}
func main() {
	env := ""
	viper.SetConfigName("appconfig")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	} else {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
		testingEnvEnable := viper.GetString("test.enable")
		devEnvEnable := viper.GetString("dev.enable")
		prodEnvEnable := viper.GetString("prod.enable")
		router := gin.New()
		//var router *gin.Engine
		if testingEnvEnable == "true" {
			env = "test"
			//router = setupRouter(true)
			gin.SetMode(gin.DebugMode)
		}
		if devEnvEnable == "true" {
			env = "dev"
			//router = setupRouter(true)
			gin.SetMode(gin.DebugMode)
		}
		if prodEnvEnable == "true" {
			env = "prod"
			//router = setupRouter(false)
			gin.SetMode(gin.ReleaseMode)
		}
		envLogLevel := fmt.Sprintf("%s.log.level", env)
		envLogFile := fmt.Sprintf("%s.log.file", env)
		envDataBaseName := fmt.Sprintf("%s.database.database", env)
		envDataBaseUser := fmt.Sprintf("%s.database.username", env)
		envDataBasePass := fmt.Sprintf("%s.database.password", env)
		envDataBaseAddress := fmt.Sprintf("%s.database.host", env)
		envServerPort := fmt.Sprintf("%s.server.port", env)
		envServerHost := fmt.Sprintf("%s.server.host", env)

		logLevel := viper.GetString(envLogLevel)
		var level model.LogLevel
		switch logLevel {
		case "debug":
			level = model.DEBUG
		case "info":
			level = model.INFO
		case "error":
			level = model.ERROR
		case "warning":
			level = model.WARNING
		case "panic":
			level = model.PANIC
		default:
			level = model.ERROR
		}
		logFile := viper.GetString(envLogFile)
		dataBaseName := viper.GetString(envDataBaseName)
		dataBaseUser := viper.GetString(envDataBaseUser)
		dataBasePass := viper.GetString(envDataBasePass)
		dataBaseHost := viper.GetString(envDataBaseAddress)
		serverPort := viper.GetString(envServerPort)
		serverHost := viper.GetString(envServerHost)
		utility.MLog.New(logFile, level)
		gin.DisableConsoleColor()
		f, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		gin.DefaultWriter = io.MultiWriter(f)
		router.Use(gin.Logger())
		v1 := router.Group("api/v1")
		{
			v1.GET("/users", GetUsers)
		}
		setupRouter(router)
		controller.Data.New(dataBaseUser, dataBasePass, dataBaseHost, dataBaseName)

		defer controller.Data.Close()
		address := fmt.Sprintf("%v:%s", serverHost, serverPort)
		srv := &http.Server{
			Addr:    address,
			Handler: router,
		}
		go func() {
			utility.MLog.Info("Server is listening port: ", address)
			if err := srv.ListenAndServe(); err != nil {
				utility.MLog.Info("Server error: ", err)
			}
		}()
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		utility.MLog.Info("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			utility.MLog.Panic("Server Shutdown:", err)
		}
		utility.MLog.Info("Server exiting")
	}
}
