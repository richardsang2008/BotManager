package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/richardsang2008/BotManager/utility"
	"github.com/spf13/viper"
	"github.com/weilunwu/go-geofence"


	"bufio"
	"encoding/json"
	"github.com/kellydunn/golang-geo"
	"github.com/richardsang2008/BotManager/controller"
	"github.com/richardsang2008/BotManager/model"
	"github.com/richardsang2008/BotManager/services"
	"io"
	"os"
	"strings"
	"sync"
)

func GetUsers(c *gin.Context) {

	c.JSON(200, "hello world")
}

const ConfigPath = "config"
const PokemonEn = "pokemon_en"

//const CpMultipliers = "cp_multipliers"
//const BaseStats = "base_stats"

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
	wg := &sync.WaitGroup{}
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


		logLevel := viper.GetString(envLogLevel)

		var level model.LogLevel
		switch logLevel {
		case "debug":
			level = model.LogLevelDEBUG
		case "info":
			level = model.LogLevelINFO
		case "error":
			level = model.LogLevelERROR
		case "warning":
			level = model.LogLevelWARNING
		case "panic":
			level = model.LogLevelPANIC
		default:
			level = model.LogLevelERROR
		}
		logFile := viper.GetString(envLogFile)

		//init the cache
		utility.MCache.New(5, 10)

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
		err:=json.Unmarshal([]byte(filterrstr), regions)
		if err != nil {
			utility.MLog.Error(err)
		}
		//create a geofence
		var genfence_zones []model.GeoFences
		if regions.Regions != nil {
			var geo_Fence = model.GeoFences{}
			for _, element := range regions.Regions {
				geo_Fence.Region = element.RegionName
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
		//load the teams data
		teamsMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "teams")
		//load the rarity data
		//rarityMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "rarity")
		//load the sizes data
		//sizesMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "sizes")
		//utility.MLog.Debug(teamsMap[2])
		//utility.MLog.Debug(pokemonMap[29])
		//utility.MLog.Debug(moveMap[135])
		//utility.MLog.Debug(rarityMap[2])
		//utility.MLog.Debug(sizesMap[3])
		//load the types data
		//typesMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "types")
		//utility.MLog.Debug(typesMap[15])
		//load the weather data
		//weatherMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "weather")
		//utility.MLog.Debug(weatherMap[6])
		//load the forms 201 data
		//forms201Map, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "forms.201")
		//utility.MLog.Debug(forms201Map[10])
		//load the forms 351 data
		//forms351Map, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "forms.351")
		//utility.MLog.Debug(forms351Map[31])
		//load the forms 386 data
		//forms386Map, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "forms.386")
		//utility.MLog.Debug(forms386Map[35])
		//load the day_or_night data
		//dayOrNightMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "day_or_night")
		//utility.MLog.Debug(dayOrNightMap[1])
		//load the leaders data
		//leadersMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "leaders")
		//utility.MLog.Debug(leadersMap[3])
		//load the severity data
		//severityMap, _ := utility.MUtility.LoadDataToIntMapString(ConfigPath, PokemonEn, "severity")
		//utility.MLog.Debug(severityMap[3])
		//load the misc data
		//miscMap, _ := utility.MUtility.LoadDataToStringMapString(ConfigPath, PokemonEn, "misc")
		//utility.MLog.Debug(miscMap["boosted"])
		//load the cp_multipliers
		//cpMultipliersMap, _ := utility.MUtility.LoadDataToFloat64MapString(ConfigPath, CpMultipliers, "cp_multipliers")
		//utility.MLog.Debug(cpMultipliersMap[7.5])
		//load the base stats data
		//baseStatsMap, _ := utility.MUtility.LoadDataToIntMapInterface(ConfigPath, BaseStats, "base_stats")
		//utility.MLog.Debug(baseStatsMap[5])

		//load the testing data from testdata.txt
		f, err1 := os.Open("test/testdata.txt")
		if err1 != nil {
			utility.MLog.Error(err1)
		}
		defer f.Close()
		var usersfilterstr []string
		f2, err2 := os.Open("test/filters.json")
		if err2 != nil {
			utility.MLog.Error(err2)
		}
		usersFiltersScanner := bufio.NewScanner(f2)
		defer f2.Close()
		for usersFiltersScanner.Scan() {
			line := usersFiltersScanner.Text()
			if !strings.HasPrefix(line, "#") || len(line) == 0 {
				usersfilterstr = append(usersfilterstr, line)
			}
		}
		messageScanner := bufio.NewScanner(f)
		for messageScanner.Scan() {
			messageAline := messageScanner.Text()
			if !strings.HasPrefix(messageAline, "#") || len(messageAline) == 0 {
				//add the string into array
				if len(usersfilterstr) > 0 {
					controller.FilterPokeMinerInputForAllUsers([]byte(messageAline), usersfilterstr, genfence_zones, pokemonMap, moveMap, teamsMap)
				}
			}
		}

		envDataBaseName := fmt.Sprintf("%s.database.database", env)
		envDataBaseUser := fmt.Sprintf("%s.database.username", env)
		envDataBasePass := fmt.Sprintf("%s.database.password", env)
		envDataBaseAddress := fmt.Sprintf("%s.database.host", env)
		dataBaseName := viper.GetString(envDataBaseName)
		dataBaseUser := viper.GetString(envDataBaseUser)
		dataBasePass := viper.GetString(envDataBasePass)
		dataBaseHost := viper.GetString(envDataBaseAddress)
		controller.Data.New(dataBaseUser, dataBasePass, dataBaseHost, dataBaseName)
		defer controller.Data.Close()
		//slack hosting
		envSlackMasterslackToken :=fmt.Sprintf("%s.slack.slackApi.masterslackToken",env)
		envSlackBotslackToken := fmt.Sprintf("%s.slack.slackApi.botslackToken",env)
		envSlackLisaslacktoken := fmt.Sprintf("%s.slack.slackApi.lisaslacktoken",env)
		envMessageQueueConsumerlookupaddress := fmt.Sprintf("%s.messagequeue.consumerlookupaddress",env)
		messageQueueConsumerlookupaddress:=viper.GetString(envMessageQueueConsumerlookupaddress)
		envMessageQueueProduceraddress := fmt.Sprintf("%s.messagequeue.produceraddress",env)
		messageQueueProduceraddress:=viper.GetString(envMessageQueueProduceraddress)
		slackMasterslackToken :=viper.GetString(envSlackMasterslackToken)
		slackBotslackToken := viper.GetString(envSlackBotslackToken)
		slackLisaslacktoken := viper.GetString(envSlackLisaslacktoken)
		slackcontroller :=controller.SlackController{}
		msgChannel:="slack_messages_channel"
		msgTopic :="slack_poke_messages"

		slackcontroller.SlackSelfHost(env,slackLisaslacktoken,slackMasterslackToken,slackBotslackToken,messageQueueProduceraddress,messageQueueConsumerlookupaddress,msgTopic,msgChannel,wg)

		//webhosting
		/*
		serverPort := viper.GetString(envServerPort)
		serverHost := viper.GetString(envServerHost)
		envServerPort := fmt.Sprintf("%s.server.port", env)
		envServerHost := fmt.Sprintf("%s.server.host", env)
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
		utility.MLog.Info("Server exiting")*/
	}
}
