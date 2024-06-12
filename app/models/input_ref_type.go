package models

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	Branch string = "branch"
	PR     string = "pr"
)

var filters []string = []string{
	"CI_COMMIT_BRANCH",
}

type RefType struct {
	Input
	Event string
}

func NewRefType() *RefType {
	return &RefType{
		Input: Input{
			Type:    Ref,
			Prefix:  "",
			Suffix:  "",
			Enabled: true,
		},
		Event: Branch,
	}
}

// Generate the final output
func (r *RefType) Generate() string {
	sanitizedBranchName, isEnvSet := os.LookupEnv("CI_COMMIT_REF_SLUG")
	if !isEnvSet {
		log.Error("CI_COMMIT_REF_SLUG not set")
		return ""
	}
	return fmt.Sprintf("%s%s%s", r.Prefix, sanitizedBranchName, r.Suffix)
}

func (r *RefType) IsEnabled() bool {
	return r.Enabled
}

func (r *RefType) Enable() {
	r.Enabled = true
}
func (r *RefType) Disable() {
	r.Enabled = false
}

func (r *RefType) SetPrefix(s string) {
	r.Prefix = s
}
func (r *RefType) SetSuffix(s string) {
	r.Suffix = s
}

func (r *RefType) SetEvent(event string) {
	r.Event = event
}

func (r *RefType) GitlabFilterPassed() bool {
	for _, filter := range filters {
		if _, exist := os.LookupEnv(filter); !exist {
			log.Warn("Gitlab filter not passed for Ref type: ", filter)
			return false
		}
	}
	return true
}

func (r *RefType) String() string {
	return fmt.Sprintf("Type: %s, Prefix: %s, Suffix: %s, Enabled: %t, Event: %s", r.Type, r.Prefix, r.Suffix, r.Enabled, r.Event)
}
