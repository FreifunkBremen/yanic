package api

import (
  "net/http"
  "encoding/json"
)

func jsonOutput(w http.ResponseWriter,data interface{}){
  js, err := json.Marshal(data)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}
