package models

import "github.com/hashicorp/go-version"

type SemverSortable []string

func (s SemverSortable) Len() int {
	return len(s)
}

func (s SemverSortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SemverSortable) Less(i, j int) bool {
	v1, err := version.NewVersion(s[i])
	if err != nil {
		panic(err)
	}
	v2, err := version.NewVersion(s[j])
	if err != nil {
		panic(err)
	}
	return v1.LessThan(v2)
}
