package helpers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func JSON(w http.ResponseWriter, status int, message string, data ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]any{
		"message": message,
	}

	if len(data) > 0 {
		response["data"] = data[0]
	}

	json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func ExtractIDFromPath(r *http.Request) string {
	pathParts := strings.Split(r.URL.Path, "/")
	return pathParts[len(pathParts)-1]
}

func ParseUUIDOr400(w http.ResponseWriter, id string) (uuid.UUID) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil
	}
	return uid
}
