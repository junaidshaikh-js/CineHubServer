package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func ParseID(idStr string) (int, error) {
	return strconv.Atoi(idStr)
}

func WriteJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	return nil
}
