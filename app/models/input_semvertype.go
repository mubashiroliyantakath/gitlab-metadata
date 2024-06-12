package models

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/mubashiroliyantakath/gitlab-metadata/app/global"
	"github.com/mubashiroliyantakath/gitlab-metadata/app/lib"
	log "github.com/sirupsen/logrus"
)

const semVerRegex = `^.*(\d+.\d+.\d+)(?:[-!$%^&*()_+|~={}\[\]:";'<>?,.\/]?([a-zA-Z]+)[-!$%^&*()_+|~={}\[\]:";'<>?,.\/]?([0-9]+))?$`

var semverPattern = regexp.MustCompile(semVerRegex)

type SemverType struct {
	Input
}

func NewSemverType() *SemverType {
	return &SemverType{
		Input: Input{
			Type:    Semver,
			Enabled: true,
		},
	}
}

func (s *SemverType) Generate() string {
	if _, isEnvSet := os.LookupEnv("CI_COMMIT_TAG"); !isEnvSet {
		log.Warn("CI_COMMIT_TAG not set. Using non-standard logic to set the version number.")
		return dynamicVersionGeneration()
	} else {
		log.Info("CI_COMMIT_TAG set. Parsing the tag and creating the version number.")
		return versionGeneration()
	}

}

func (s *SemverType) IsEnabled() bool {
	return s.Enabled
}

func (s *SemverType) Enable() {
	s.Enabled = true

}
func (s *SemverType) Disable() {
	s.Enabled = false
}

func (s *SemverType) SetPrefix(string) {
	log.Debug("SemverType does not have a prefix. Disregarding.")

}
func (s *SemverType) SetSuffix(string) {
	log.Debug("SemverType does not have a suffix. Disregarding.")
}

func (s *SemverType) SetEvent(string) {
	log.Debug("SemverType does not have an event. Disregarding.")
}

func (r *SemverType) GitlabFilterPassed() bool {
	return true
}

func (s *SemverType) String() string {
	return fmt.Sprintf("Type: %s, Prefix: %s, Suffix: %s, Enabled: %t", s.Type, s.Prefix, s.Suffix, s.Enabled)
}

func dynamicVersionGeneration() string {
	_, iidExists := os.LookupEnv("CI_PIPELINE_IID")
	if !iidExists {
		log.Fatal("CI_PIPELINE_IID not set. Not a Gitlab Pipeline. Exiting.")
	}
	gitFetch := exec.Command("git", "fetch", "--tags", "-f")
	err := gitFetch.Run()
	if err != nil {
		log.Error("Failed to fetch the tags: ", err)
		return ""
	}
	getTags := exec.Command("git", "tag", "--sort=creatordate")
	output, err := getTags.Output()
	if err != nil {
		log.Error("Failed to fetch the latest tag: ", err)
		return ""
	}
	tagList := strings.Split(string(output), "\n")
	tagList = lib.DeleteEmptyFromList(tagList)
	if len(tagList) == 0 {
		log.Warn("No tags found. Using a default version number.")
		tagList = append(tagList, "0.1.0")
	}
	log.Debug("All available tags list: ", tagList)
	var semverTagList SemverSortable
	for _, tag := range tagList {
		if semverPattern.MatchString(tag) {
			matches := semverPattern.FindStringSubmatch(tag)
			if matches[2] != "" {
				log.Debug("Skipping pre-release tag: ", tag)
			} else {
				sanitizedTag := lib.SanitizeTags(matches[1])
				semverTagList = append(semverTagList, sanitizedTag)
			}

		} else {
			log.Debug("Skipping non-semver tag: ", tag)
		}
	}

	log.Debug("Semver tags list: ", semverTagList)
	sort.Sort(semverTagList)
	log.Debug("Sorted semver tags list: ", semverTagList)
	latestTag := semverTagList[len(semverTagList)-1]
	log.Debug("Latest tag: ", latestTag)
	latestTag = fmt.Sprintf("%v.%v", latestTag, os.Getenv("CI_PIPELINE_IID"))
	log.Debug("Generated version: ", latestTag)
	log.Info("Version Generated: ", latestTag)
	global.Version = latestTag
	return latestTag

}

func versionGeneration() string {
	tag := os.Getenv("CI_COMMIT_TAG")
	if semverPattern.MatchString(tag) {
		matches := semverPattern.FindStringSubmatch(tag)
		if matches[2] != "" {
			log.Info("Found a pre-release tag: ", tag)
			version := matches[1]
			preReleaseKind := matches[2]
			preReleaseCounter := matches[3]
			tag := lib.SanitizePreReleaseTags(lib.SanitizeTags(version), preReleaseKind, preReleaseCounter)
			log.Info("Version Generated: ", tag)
			global.Version = tag
			return tag

		} else {
			sanitizedTag := lib.SanitizeTags(matches[1])
			log.Debug("Sanitized tag: ", sanitizedTag)
			log.Info("Version Generated: ", sanitizedTag)
			global.Version = sanitizedTag
			return sanitizedTag
		}
	}
	return ""
}
