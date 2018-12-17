package api

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/euforia/thrap/pkg/thrap"
	nomad "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/client/driver/env"
	"github.com/hashicorp/nomad/jobspec"
)

// This file contains backward compatibility code which should eventually
// be deprecated

var (
	metaPlaceHolderRe = regexp.MustCompile(`\<[a-zA-Z0-9_\-\.]+\>`)
)

const (
	secretTmplSMarker = `{{ with printf `
	secretTmplEMarker = `|`
	secretTmplReplace = `"%s" (env "NOMAD_META_` + thrap.SecretsPathVarName + `") `
)

// replaceMetaPlaceholders replaces internal placeholders with nomad variables
// This will eventually be deprecated
func replaceMetaPlaceholders(data string) string {
	return metaPlaceHolderRe.ReplaceAllStringFunc(data, func(match string) string {
		// Remove leading < and trailing >
		stripped := match[1 : len(match)-1]
		// Append nomad meta and the mapped value
		return fmt.Sprintf("${%s%s}", env.MetaPrefix, stripped)
	})
}

// // replaceMetaPlaceholders replaces internal placeholders with nomad variables
// // This will eventually be deprecated
// func replaceMetaPlaceholders(data string, maps ...map[string]string) string {
// 	return metaPlaceHolderRe.ReplaceAllStringFunc(data, func(match string) string {
// 		// Remove leading < and trailing >
// 		stripped := match[1 : len(match)-1]
// 		for _, m := range maps {
// 			if val, ok := m[stripped]; ok {
// 				// Append nomad meta and the mapped value
// 				return fmt.Sprintf("${%s%s}", env.MetaPrefix, val)
// 			}
// 		}
// 		return match
// 	})
// }

func parseMoldHCLDeployDesc(in []byte) (*pb.DeploymentDescriptor, error) {
	// Replace < > first
	str := replaceMetaPlaceholders(string(in))

	job, err := jobspec.Parse(bytes.NewBuffer([]byte(str)))
	if err != nil {
		return nil, err
	}

	// Remove meta.env_type
	remove := "${meta.env_type}"
	csts := make([]*nomad.Constraint, 0, len(job.Constraints))
	for _, cst := range job.Constraints {
		if cst.LTarget != remove {
			csts = append(csts, cst)
		}
	}
	job.Constraints = csts

	removeVarsFromMetaVals(job.Meta)

	for _, group := range job.TaskGroups {
		removeVarsFromMetaVals(group.Meta)

		for _, task := range group.Tasks {
			removeVarsFromMetaVals(task.Meta)

			for _, tmpl := range task.Templates {
				data := tmpl.EmbeddedTmpl
				if data == nil {
					continue
				}
				// Replace secrets path with the variablized one
				replaced := replaceSecretsTemplateVars(*data)
				tmpl.EmbeddedTmpl = &replaced
			}
		}
	}

	return makeNomadJSONDeployDesc(job)
}

// empties out values in the meta block that contain variables
func removeVarsFromMetaVals(meta map[string]string) {
	for k, v := range meta {
		if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
			meta[k] = ""
		}
	}
}

// replaces the secrets key
func replaceSecretsTemplateVars(s string) string {
	i := strings.Index(s, secretTmplSMarker)
	if i < 0 {
		return ""
	}
	i += 15

	e := strings.Index(s[i:], secretTmplEMarker)
	if e < 1 {
		return ""
	}
	e += i

	pre := s[:i]

	post := s[e:]

	return pre + secretTmplReplace + post
}
