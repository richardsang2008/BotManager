package main
import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/richardsang2008/BotManager/utility"

)

func GetUsers(c *gin.Context) {
//	log.Debug("I am here")
	c.JSON(200, "hello world")
}

type User struct {
	gorm.Model
	ID        uint   `gorm:"primary_key`
	Uname     string `sql:"type:VARCHAR(255)"`
	CreatedAt time.Time
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
		if testingEnvEnable == "true" {
			env = "test"
		}
		if devEnvEnable == "true" {
			env = "dev"
		}
		if prodEnvEnable == "true" {
			env = "prod"
		}
		envLogLevel := fmt.Sprintf("%s.log.level", env)
		envLogFile := fmt.Sprintf("%s.log.file", env)
		envDataBaseName := fmt.Sprintf("%s.database.database", env)
		envDataBaseUser := fmt.Sprintf("%s.database.username", env)
		envDataBasePass := fmt.Sprintf("%s.database.password", env)
		envServerPort := fmt.Sprintf("%s.server.port", env)
		logLevel := viper.GetString(envLogLevel)
		logFile := viper.GetString(envLogFile)
		dataBaseName := viper.GetString(envDataBaseName)
		dataBaseUser := viper.GetString(envDataBaseUser)
		dataBasePass := viper.GetString(envDataBasePass)
		serverPort := viper.GetString(envServerPort)
		log := utility.Log{}
		Log.NewLogger(logFile, logLevel)

		router := gin.Default()
		v1 := router.Group("api/v1")

		{
			v1.GET("/users", GetUsers)
		}
		dbConnectionStr := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8&parseTime=True&loc=Local", dataBaseUser, dataBasePass, dataBaseName)
		db, err := gorm.Open("mysql", dbConnectionStr)
		db.CreateTable(&User{})
		defer db.Close()
		if err != nil {
			log.Panic("DB is not open ")
		}

		ports := fmt.Sprintf(":%s", serverPort)
		log.Info("Server is running from port ", serverPort)
		router.Run(ports)
	}

}
