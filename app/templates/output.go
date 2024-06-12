package templates

// Template for output file
const OutputEnvFileTemplate = `VERSION={{- .Version }}
TAGS={{- range $index, $element := .Tags -}}
{{- if $index }},{{ end -}}
{{- $element -}}
{{- end -}}
`
