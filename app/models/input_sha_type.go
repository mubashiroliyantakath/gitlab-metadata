package models

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type ShaType struct {
	Input
}

func NewShaType() *ShaType {
	return &ShaType{
		Input: Input{
			Type:    Sha,
			Enabled: true,
		},
	}
}

func (s *ShaType) Generate() string {
	shortSha, isEnvSet := os.LookupEnv("CI_COMMIT_SHORT_SHA")
	if !isEnvSet {
		log.Error("CI_COMMIT_SHORT_SHA not set")
		return ""
	} else {
		return fmt.Sprintf("%s%s%s", s.Prefix, shortSha, s.Suffix)
	}

}

func (s *ShaType) IsEnabled() bool {
	return s.Enabled
}

func (s *ShaType) Enable() {
	s.Enabled = true

}
func (s *ShaType) Disable() {
	s.Enabled = false
}

func (s *ShaType) SetPrefix(p string) {
	s.Prefix = p

}
func (s *ShaType) SetSuffix(p string) {
	s.Suffix = p
}

func (s *ShaType) SetEvent(string) {
	log.Debug("ShaType does not have an event. Disregarding.")
}

func (r *ShaType) GitlabFilterPassed() bool {
	log.Debug("No filtering required for sha type")
	return true
}

func (s *ShaType) String() string {
	return fmt.Sprintf("Type: %s, Prefix: %s, Suffix: %s, Enabled: %t", s.Type, s.Prefix, s.Suffix, s.Enabled)
}
