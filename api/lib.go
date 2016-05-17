package api

import (
	"bytes"
	"strings"
	"net/http"
	"encoding/json"
	"encoding/base64"

	"github.com/julienschmidt/httprouter"
)

func jsonOutput(w http.ResponseWriter, r *http.Request,data interface{}){
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers","Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write(js)
}
func BasicAuth(h httprouter.Handle, pass []byte) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			if origin := r.Header.Get("Origin"); origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers","Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			const basicAuthPrefix string = "Basic "

			// Get the Basic Authentication credentials
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, basicAuthPrefix) {
					// Check credentials
					payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
					if err == nil {
							pair := bytes.SplitN(payload, []byte(":"), 2)
							if len(pair) == 2 &&
									bytes.Equal(pair[1], pass) {

									// Delegate request to the given handle
									h(w, r, ps)
									return
							}
					}
			}

			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
}
