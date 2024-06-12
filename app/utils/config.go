package utils

import (
	"bytes"
	"embed"
	"os"
	"strings"

	"github.com/mubashiroliyantakath/gitlab-metadata/app/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	AppConfig *models.MetadataInput
)

//go:embed conf.yaml
var configFile embed.FS

func NewAppConfig() error {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("MA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if _, err := os.Stat("conf.yaml"); os.IsNotExist(err) {
		log.Debug("conf.yaml does not exist in the root path. Proceeding with the embedded config.")
		configFileContent, err := configFile.ReadFile("conf.yaml")
		if err != nil {
			log.Fatal("Failed to read embedded config file: ", err)
		}
		//
		err = v.ReadConfig(bytes.NewBuffer(configFileContent))
		if err != nil {
			return err
		}
	} else {
		v.SetConfigFile("conf.yaml")
		err := v.ReadInConfig()
		if err != nil {
			return err
		}
	}

	err := v.Unmarshal(&AppConfig)
	if err != nil {
		return err
	}
	log.Debug("AppConfig loaded: ", AppConfig)
	return nil
}
