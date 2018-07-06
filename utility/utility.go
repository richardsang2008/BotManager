package utility

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"strconv"
)

var (
	MLog         Log
	MCache       Cache
	MUtility     Utility
	MPokeUtility PokeUtility
	MSlackUtility SlackUtility

)

type Utility struct {
}

func (u *Utility) LoadDataToStringMapString(configPath string, filename string, stringname string) (map[string]string, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		MLog.Panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return nil, err
	} else {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			MLog.Info("Config file changed:", e.Name)
		})
		datamap := v.GetStringMapString(stringname)
		return datamap, nil
	}
}
func (u *Utility) LoadDataToFloat64MapString(configPath string, filename string, stringname string) (map[float64]string, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		MLog.Panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return nil, err
	} else {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			MLog.Info("Config file changed:", e.Name)
		})
		datamap := v.GetStringMapString(stringname)
		//change datamap into map[int]string
		intmap := make(map[float64]string)
		for k, v := range datamap {
			key, _ := strconv.ParseFloat(k, 64)
			intmap[key] = v
		}
		return intmap, nil
	}
}
func (u *Utility) LoadDataToIntMapString(configPath string, filename string, stringname string) (map[int]string, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		MLog.Panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return nil, err
	} else {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			MLog.Info("Config file changed:", e.Name)
		})
		datamap := v.GetStringMapString(stringname)
		//change datamap into map[int]string
		intmap := make(map[int]string)
		for k, v := range datamap {
			key, _ := strconv.Atoi(k)
			intmap[key] = v
		}
		return intmap, nil
	}
}
func (u *Utility) LoadDataToIntMapInterface(configPath string, filename string, stringname string) (map[int]interface{}, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		MLog.Panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return nil, err
	} else {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			MLog.Info("Config file changed:", e.Name)
		})
		dmap := v.GetStringMap(stringname)
		x := make(map[int]interface{})
		for k, v := range dmap {
			intk, _ := strconv.Atoi(k)
			x[intk] = v
		}
		return x, nil
	}
}

// FileToMap is a helper function that combines ReadFile and ToMap()
func (u *Utility) FileToMap(filepath string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	dataMap, err := u.ToMap(data)
	return dataMap, err
}

// ToMap assumes top level keys are strings (i.e. NOT a json file with just []) and returns the bytes as a Map of objects
func (u *Utility) ToMap(data []byte) (map[string]interface{}, error) {
	// A map of string to any type https://blog.golang.org/laws-of-reflection , http://research.swtch.com/interfaces
	var datamap map[string]interface{}
	err := json.Unmarshal(data, &datamap)
	return datamap, err
}
