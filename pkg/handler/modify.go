package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var defaultFromParam = "file"
var defaultFileParam = "file"

type jsonFormatedImage struct {
	Base64 string `json:"file"`
}

// ModifyEndpoint just illustrates work with service. There is no need
// to create file on each request. Certain actions depend on business requirements.
func (h *handler) ModifyEndpoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := h.getFile(r)
		if err != nil {
			errHTTP(w, err)
			return
		}

		ID, size, err := h.service.Create(file)
		if err != nil {
			errHTTP(w, err)
			return
		}
		reader, err := h.service.Open(ID)
		if err != nil {
			errHTTP(w, err)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
		mods, err := h.extractor.Extract(r)

		if err != nil {
			errHTTP(w, err)
			return
		}

		h.service.Modify(reader, w, mods)
	}
}

// getFile returns io.Reader of file provided in request.
// Current API accept files in three ways:
// - GET: URL in GET-param. Name of the param is "file" by default
// - POST: multipart/form-data. Name of the param is "file" by default
// - POST: as base64-encoded string within "file" property of JSON object
func (h *handler) getFile(r *http.Request) (io.Reader, error) {
	switch r.Method {
	case http.MethodGet:
		from := r.URL.Query().Get(defaultFromParam)
		file, err := http.Get(from)
		if err == nil {
			return file.Body, nil
		}
	case http.MethodPost:
		file, _, err := r.FormFile(defaultFileParam)
		if err != nil {
			// multipart/form-data failed, tring parse json
			decoder := json.NewDecoder(r.Body)
			var j jsonFormatedImage
			err := decoder.Decode(&j)
			if err != nil {
				return nil, &BadRequestError{err}
			}

			// The actual image data starts after the ","
			i := strings.Index(j.Base64, ",")
			if i < 0 {
				return nil, &BadRequestError{fmt.Errorf("Malformed base64, no comma")}
			}

			// decode base64 to image
			return base64.NewDecoder(base64.StdEncoding, strings.NewReader(j.Base64[i+1:])), nil
		}
		return file, err
	}

	return nil, &BadRequestError{fmt.Errorf("Image file not found in request")}
}
