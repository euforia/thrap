package api

import (
	"encoding/json"
	"net/http"
)

// func writeJSON(w http.ResponseWriter, respData interface{}) {
// 	data, err := json.Marshal(respData)
// 	if err != nil {
// 		w.WriteHeader(500)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	setAccessControl(w)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(data)
// }

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
