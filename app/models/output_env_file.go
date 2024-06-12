package models

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mubashiroliyantakath/gitlab-metadata/app/global"
	"github.com/mubashiroliyantakath/gitlab-metadata/app/templates"
)

type OutputEnvFile struct {
	Output
}

type OutputData struct {
	Version string
	Tags    []string
}

func NewOutputEnvFile() *OutputEnvFile {
	tpl, err := template.New("env-file").Parse(templates.OutputEnvFileTemplate)
	if err != nil {
		log.Fatal("Failed to parse template: ", err)
	}
	return &OutputEnvFile{
		Output: Output{
			Template: tpl,
			FileName: filepath.Join(global.CiProjectDir, "metadata-out.env"),
			Data: &OutputData{
				Version: global.Version,
				Tags:    global.Tags,
			},
		},
	}
}

func (o *OutputEnvFile) Write() error {
	outputFile, err := os.Create(o.FileName)
	if err != nil {
		log.Fatal("Failed to create output file: ", err)
	}
	defer outputFile.Close()
	err = o.Template.Execute(outputFile, o.Data)
	if err != nil {
		log.Fatal("Failed to write output file: ", err.Error())
	}
	return nil
}
