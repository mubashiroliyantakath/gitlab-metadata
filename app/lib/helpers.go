package lib

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

func DeleteEmptyFromList(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func SanitizeTags(rawTag string) string {
	sanitizedTag := strings.ReplaceAll(rawTag, "_", ".")
	sanitizedTag = strings.ReplaceAll(sanitizedTag, "-", ".")
	log.Debug("Sanitized tag: ", sanitizedTag)
	return sanitizedTag
}

func SanitizePreReleaseTags(sanitizedVersion string, preReleaseKind string, preReleaseCounter string) string {
	log.Debug("Pre-release kind: ", preReleaseKind)
	log.Debug("Pre-release counter: ", preReleaseCounter)
	return fmt.Sprintf("%s-%s%s", sanitizedVersion, preReleaseKind, preReleaseCounter)
}
