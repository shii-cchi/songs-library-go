package delivery

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON sends a JSON response with the specified HTTP status code and payload.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, _ := json.Marshal(payload)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
