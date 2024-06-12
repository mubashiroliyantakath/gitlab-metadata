package utils

import (
	log "github.com/sirupsen/logrus"

	"github.com/mubashiroliyantakath/gitlab-metadata/app/interfaces"
)

func GenerateTags(inputs []interfaces.Input) []string {
	var tags []string
	for _, item := range inputs {
		if item.IsEnabled() {
			tags = append(tags, item.Generate())
		}
	}
	if len(tags) == 0 {
		log.Fatal("No tags were generated")
	}
	return tags
}
