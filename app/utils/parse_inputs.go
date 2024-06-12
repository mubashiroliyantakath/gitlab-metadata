package utils

import (
	"os"
	"strconv"
	"strings"

	"github.com/mubashiroliyantakath/gitlab-metadata/app/interfaces"
	"github.com/mubashiroliyantakath/gitlab-metadata/app/models"

	log "github.com/sirupsen/logrus"
)

func ParseInputs(inputs []string) []interfaces.Input {

	var tagList []interfaces.Input

	if len(inputs) == 0 {
		log.Warn("No inputs were provided. Will proceed with default inputs.")
		if _, exist := os.LookupEnv("CI_COMMIT_BRANCH"); exist {
			log.Info("Branch input detected")
			tagList = append(tagList, models.NewRefType())
		}

		tagList = append(tagList, models.NewSemverType())

		tagList = append(tagList, models.NewShaType())
		for number, tag := range tagList {
			log.Debug("Default input ", number+1, " loaded => ", tag)
		}
		return filterInputs(tagList)
	}

	for _, input := range inputs {
		log.Info("Processing input: ", input)
		var inputItem interfaces.Input
		var enabled bool
		var suffix string
		var prefix string

		inputMap := make(map[string]string)
		pairs := strings.Split(input, ",")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				inputMap[kv[0]] = kv[1]
			}
		}

		if _, ok := inputMap["type"]; !ok {
			log.Error("Input type is required")
			continue
		}

		if _, ok := inputMap["enabled"]; ok {
			enabled, _ = strconv.ParseBool(inputMap["enabled"])
		} else {
			enabled = true
		}

		if _, ok := inputMap["suffix"]; ok {
			suffix = inputMap["suffix"]
		} else {
			suffix = ""
		}

		if _, ok := inputMap["prefix"]; ok {
			prefix = inputMap["prefix"]
		} else {
			prefix = ""
		}

		// Determine the input type
		switch inputMap["type"] {
		// Workflow when type is "ref"
		case models.Ref:
			if _, ok := inputMap["event"]; !ok {
				log.Error("Input event is required")
				log.Warn("The following input is disregarded: ", input)
				continue
			}
			inputItem = models.NewRefType()
			switch inputMap["event"] {
			case models.Branch:
				log.Info("Tag with branch")
				inputItem.SetEvent(models.Branch)
			case models.PR:
				log.Info("Tag with PR")
				inputItem.SetEvent(models.PR)
			default:
				log.Error("Invalid ref event")
				log.Warn("The following input is disregarded: ", input)
				continue
			}

		case models.Semver:
			log.Info("Tag with semver")
			inputItem = models.NewSemverType()
		case models.Sha:
			log.Info("Tag with sha")
			inputItem = models.NewShaType()
		default:
			log.Error("Invalid input type")
			log.Warn("The following input is disregarded: ", &input)
		}

		if prefix != "" {
			inputItem.SetPrefix(prefix)
		}
		if suffix != "" {
			inputItem.SetSuffix(suffix)
		}
		if enabled {
			inputItem.Enable()
		} else {
			inputItem.Disable()
		}
		log.Info("Input processed successfully")
		log.Info(inputItem.String())
		tagList = append(tagList, inputItem)
	}
	if len(tagList) == 0 {
		log.Fatal("No valid inputs were found")
	} else {
		log.Debug("Inputs parsed successfully")
		for number, input := range tagList {
			log.Debug("User supplied input ", number+1, " loaded => ", input)
		}
		return filterInputs(tagList)
	}
	return filterInputs(tagList)
}
