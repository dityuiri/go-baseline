package util

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
)

func HttpResponseBodyParser(response *http.Response, data interface{}) error {
	defer response.Body.Close()
	var body = response.Body

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		var err error

		body, err = gzip.NewReader(body)
		if err != nil {
			return err
		}
	}

	// This will strip away any additional information like encoding.
	mediatype, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if body, err := io.ReadAll(body); err != nil {
		return err
	} else if mediatype == "application/json" {
		return json.Unmarshal(body, data)
	} else {
		return fmt.Errorf(string(body))
	}
}

func HttpRequestBodyParser(req *http.Request, data interface{}) error {
	defer req.Body.Close()
	var body = req.Body

	switch req.Header.Get("Content-Encoding") {
	case "gzip":
		var err error

		body, err = gzip.NewReader(body)
		if err != nil {
			return err
		}
	}

	// This will strip away any additional information like encoding.
	mediatype, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if body, err := io.ReadAll(body); err != nil {
		return err
	} else if mediatype == "application/json" {
		return json.Unmarshal(body, data)
	} else {
		return fmt.Errorf(string(body))
	}
}

func WriteResponse(w http.ResponseWriter, body interface{}, status int) {
	resp, _ := json.Marshal(body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}
