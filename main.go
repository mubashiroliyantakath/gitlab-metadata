package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/mubashiroliyantakath/gitlab-metadata/app/global"
	"github.com/mubashiroliyantakath/gitlab-metadata/app/lib"
	"github.com/mubashiroliyantakath/gitlab-metadata/app/models"
	"github.com/mubashiroliyantakath/gitlab-metadata/app/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
	log.Debug("Log level set to ", logLevel)

	gitPath, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Git is not installed or not found in PATH: ", err)
	}
	log.Debug("Git path: ", gitPath)

	err = utils.NewAppConfig()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	projectDir, envExists := os.LookupEnv("CI_PROJECT_DIR")
	if !envExists {
		log.Warn("CI_PROJECT_DIR not set. Using current directory as project directory.")
	} else {
		log.Info("Project directory set to: ", projectDir)
		global.CiProjectDir = projectDir
	}

}

func main() {
	inputs := utils.ParseInputs(strings.Fields(utils.AppConfig.Input))
	tags := utils.GenerateTags(inputs)
	log.Info("Tags generated: ", tags)
	tags = lib.DeleteEmptyFromList(tags)
	global.Tags = tags
	log.Info("tags: ", global.Tags)
	outputFile := models.NewOutputEnvFile()
	err := outputFile.Write()
	if err != nil {
		log.Fatal("Failed to write output file: ", err)
	}
}
