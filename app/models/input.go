package models

type InputType string

const (
	Ref    string = "ref"
	Semver string = "semver"
	Sha    string = "sha"
)

type Input struct {
	Type    string
	Prefix  string
	Suffix  string
	Enabled bool
}
