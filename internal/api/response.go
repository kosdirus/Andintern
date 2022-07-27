package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type RangeItemsResponse struct {
	Cars interface{} `json:"cars"`
}

// Response describes http response for api v3.
type Response struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	w.Header().Set(HeaderContentType, MimeApplicationJSON)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		log.Println("error while encoding data to respond with json:", data)
	}
}

// RespondData responds with custom status code and JSON in format: {"data": <val>}.
func RespondData(ctx context.Context, w http.ResponseWriter, code int, val interface{}) {
	RespondJSON(ctx, w, code, &Response{
		Data: val,
	})
}

// RespondDataOK responds with 200 status code and JSON in format: {"data": <val>}.
func RespondDataOK(ctx context.Context, w http.ResponseWriter, val interface{}) {
	RespondData(ctx, w, http.StatusOK, val)
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error) {
	var code = http.StatusBadRequest
	data := &Response{
		Error: err.Error(),
	}
	RespondJSON(ctx, w, code, data)
}
