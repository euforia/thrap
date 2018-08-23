package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/nomad/client/driver/env"
)

func setAccessControlHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func writeJSONResponse(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	if resp == nil {
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	setAccessControlHeaders(w)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

var (
	metaPlaceHolderRe = regexp.MustCompile(`\<[a-zA-Z0-9_\-\.]+\>`)
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
