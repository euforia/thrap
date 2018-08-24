package api

import (
	"fmt"
	"testing"
)

var testTemplate1 = `{{ with printf "secret/%s" (env "NOMAD_JOB_NAME") | secret }}{{ range $k, $v := .Data }}{{ $k }}={{ $v }}
{{ end }}
        DB_NAME = "{{.Data.DB_NAME}}"
{{ end }}`

var testTemplate2 = `{{ with printf "secret/atlas" | secret }}{{ range $k, $v := .Data }}{{ $k }}={{ $v }}
{{ end }}
        DB_NAME = "{{.Data.DB_NAME}}"
{{ end }}`

func Test_replaceSecretsTemplateVars(t *testing.T) {
	out := replaceSecretsTemplateVars(testTemplate1)
	fmt.Println(out)
	out = replaceSecretsTemplateVars(testTemplate2)
	fmt.Println(out)
}
