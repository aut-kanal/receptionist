package configuration

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configFilePath = ""
	receptionistConfig  *viper.Viper

	once sync.Once
)

// GetInstance returns an instance of viper config
func GetInstance() *viper.Viper {
	once.Do(func() {
		loadConfig()
	})
	return receptionistConfig
}

func loadConfig() {
	config := viper.New()

	// Setting defaults for this application
	config.SetDefault("addr", ":8000")
	config.SetDefault("debug", true)

	config.SetDefault("db.dialect", "sqlite3")
	config.SetDefault("db.path", "sqlite.db")
	config.SetDefault("db.host", "")
	config.SetDefault("db.port", "")
	config.SetDefault("db.name", "")
	config.SetDefault("db.user", "")
	config.SetDefault("db.password", "")

	if configFilePath != "" {
		config.SetConfigFile(configFilePath)

		config.OnConfigChange(OnConfigChanged)
		config.WatchConfig()

		err := config.ReadInConfig()
		if err != nil {
			logrus.Errorf("can't read config file, %v", err)receptionist
			receptionistConfig = config
			return
		}
		logrus.Infof("configuration file is loaded from %s", configFilePath)receptionistreceptionist
	}

	logrus.Debugf("loaded config: %v", config.AllSettings())receptionistreceptionist
	receptionistConfig = config
}

// SetFilePath sets path of config file
func SetFilePath(filePath string) {
	configFilePath = filePath
	receptionistConfig = nil
}

// SetDebugLogLevel sets log level to debug mode
func SetDebugLogLevel(isDebug bool) {
	if isDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("log level is set to Debug")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// OnConfigChanged excuates when config changes
func OnConfigChanged(_ fsnotify.Event) {
	loadConfig()
	logrus.Info("configuration is reloaded")
}
