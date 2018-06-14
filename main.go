package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/richardsang2008/BotManager/utility"
	"github.com/spf13/viper"
	"github.com/weilunwu/go-geofence"
	"time"

	"bufio"
	"context"
	"encoding/json"
	"github.com/kellydunn/golang-geo"
	"github.com/richardsang2008/BotManager/controller"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/services"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func GetUsers(c *gin.Context) {

	c.JSON(200, "hello world")
}

const ConfigPath = "config"
const PokemonEn = "pokemon_en"
const CpMultipliers = "cp_multipliers"
const BaseStats = "base_stats"

func setupRouter(router *gin.Engine) {
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
	viper.AddConfigPath(ConfigPath)
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

		//init the cache
		utility.MCache.New(5, 10)
		serverPort := viper.GetString(envServerPort)
		serverHost := viper.GetString(envServerHost)
		//load the regions inform
		fr, errr := os.Open("config/regions.json")
		if errr != nil {
			utility.MLog.Error(errr)
		}
		defer fr.Close()
		scannerr := bufio.NewScanner(fr)
		filterrstr := ""
		for scannerr.Scan() {
			line := scannerr.Text()
			if !strings.HasPrefix(line, "#") || len(line) == 0 {
				filterrstr = line
			}
		}
		regions := &model.Regions{}
		json.Unmarshal([]byte(filterrstr), regions)
		//create a geofence

		var genfence_zones []model.GeoFences

		if regions.Regions != nil {
			var geo_Fence = model.GeoFences{}
			for _, element := range regions.Regions {
				geo_Fence.Region = element.Region
				polygon := []*geo.Point{}
				for _, zone := range element.Zone {
					d := geo.NewPoint(zone.Latitude, zone.Longitude)
					polygon = append(polygon, d)

				}
				x_geofence := geofence.NewGeofence([][]*geo.Point{polygon, []*geo.Point{}})
				geo_Fence.Geofence = x_geofence
				genfence_zones = append(genfence_zones, geo_Fence)
			}
		}
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
		//load the pokemon data
		pokemonMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "pokemon")
		//load the move data
		moveMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "moves")
		//load the rarity data
		rarityMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "rarity")
		//load the sizes data
		sizesMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "sizes")
		//load the teams data
		teamsMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "teams")
		//load the types data
		typesMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "types")
		//load the weather data
		weatherMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "weather")
		//load the forms 201 data
		forms201Map, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "forms.201")
		//load the forms 351 data
		forms351Map, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "forms.351")
		//load the forms 386 data
		forms386Map, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "forms.386")
		//load the day_or_night data
		dayOrNightMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "day_or_night")
		//load the leaders data
		leadersMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "leaders")
		//load the severity data
		severityMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "severity")
		//load the misc data
		miscMap, _ := utility.MUtility.LoadDataToStringMapString(ConfigPath, PokemonEn, "misc")
		//load the cp_multipliers
		//cpMultipliersMap, _ := utility.MUtility.LoadDataToFloat64MapString(ConfigPath, CpMultipliers, "cp_multipliers")
		//load the base stats data
		//baseStatsMap, _ := utility.MUtility.LoadDataToIntMapInterface(ConfigPath, BaseStats, "base_stats")

		utility.MLog.Debug(pokemonMap[29])
		utility.MLog.Debug(moveMap[135])
		utility.MLog.Debug(rarityMap[2])
		utility.MLog.Debug(sizesMap[3])
		utility.MLog.Debug(teamsMap[2])
		utility.MLog.Debug(typesMap[15])
		utility.MLog.Debug(weatherMap[6])
		utility.MLog.Debug(forms201Map[10])
		utility.MLog.Debug(forms351Map[31])
		utility.MLog.Debug(forms386Map[35])
		utility.MLog.Debug(dayOrNightMap[1])
		utility.MLog.Debug(leadersMap[3])
		utility.MLog.Debug(severityMap[3])
		utility.MLog.Debug(miscMap["boosted"])
		//utility.MLog.Debug(cpMultipliersMap[7.5])
		//utility.MLog.Debug(baseStatsMap[5])
		//get local time
		//loc, _ := time.LoadLocation("America/Los_Angeles")
		//t := time.Now().In(loc)
		//utility.MLog.Info(t)

		//load the testing data from testdata.txt
		f, err1 := os.Open("test/testdata.txt")
		if err1 != nil {
			utility.MLog.Error(err1)
		}
		defer f.Close()
		f2, err2 := os.Open("test/filters.json")
		if err2 != nil {
			utility.MLog.Error(err2)
		}
		scanner2 := bufio.NewScanner(f2)
		defer f2.Close()
		filterstr := ""
		for scanner2.Scan() {
			line := scanner2.Text()
			if !strings.HasPrefix(line, "#") || len(line) == 0 {
				filterstr = line
			}
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			aline := scanner.Text()
			if !strings.HasPrefix(aline, "#") || len(aline) == 0 {
				general, isNotifyToUser, _ := controller.FilterPokeMinerInput([]byte(aline), []byte(filterstr), genfence_zones, pokemonMap, moveMap, teamsMap)
				if isNotifyToUser {
					utility.MLog.Debug("Send to user message")

					if err != nil {
						utility.MLog.Error(err)
					} else {
						switch t := general.(type) {
						case *model.PokeMinerMonMessage:
							mon := general.(*model.PokeMinerMonMessage)
							utility.MLog.Debug("I am pokemessage  %f , %.2f", *(mon.Message.PokemonID), *(mon.Message.Iv))
						case *model.PokeMinerRaidMessage:
							mon := general.(*model.PokeMinerRaidMessage)
							if mon.Message.Cp == nil || *(mon.Message.Cp) == 0 {
								utility.MLog.Debug("I am eggmessage", mon.Message.Latitude)
							} else {
								utility.MLog.Debug("I am raidmessage", mon.Message.Latitude)
							}

						case *model.PokeMinerGymMessage:
							mon := general.(*model.PokeMinerGymMessage)
							utility.MLog.Debug("I am gymmessage", mon.Message.Latitude)
							if mon.Message.Guards == nil {
								utility.MLog.Debug("there is no guard data")
							} else {
								utility.MLog.Debug("guard 0 id ", mon.Message.Guards[0].PokemonID)
							}
						default:
							_ = t
							utility.MLog.Debug("nothing")
						}
					}
				}
			}
		}
		/*envDataBaseName := fmt.Sprintf("%s.database.database", env)
		envDataBaseUser := fmt.Sprintf("%s.database.username", env)
		envDataBasePass := fmt.Sprintf("%s.database.password", env)
		envDataBaseAddress := fmt.Sprintf("%s.database.host", env)
		dataBaseName := viper.GetString(envDataBaseName)
		dataBaseUser := viper.GetString(envDataBaseUser)
		dataBasePass := viper.GetString(envDataBasePass)
		dataBaseHost := viper.GetString(envDataBaseAddress)
		controller.Data.New(dataBaseUser, dataBasePass, dataBaseHost, dataBaseName)
		defer controller.Data.Close()*/
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
