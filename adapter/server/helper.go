package server

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
)

type Error string

const (
	ErrUnsupportedMediaType Error = "Unsupported Media Type"
)

func (e Error) Error() string {
	return string(e)
}

func Parse(request *http.Request, data interface{}) error {
	defer request.Body.Close()

	body := request.Body

	switch request.Header.Get("Content-Encoding") {
	case "gzip":
		var err error

		body, err = gzip.NewReader(body)
		if err != nil {
			return err
		}
	}

	// This will strip away any additional information like encoding.
	mediatype, _, _ := mime.ParseMediaType(request.Header.Get("Content-Type"))
	switch mediatype {
	case "application/json":
		return json.NewDecoder(request.Body).Decode(data)
		//TODO add different mediatype unmarshaler here
	default:
		if _, err := ioutil.ReadAll(body); err != nil {
			return err
		}
		return ErrUnsupportedMediaType
	}
}

func Response(w http.ResponseWriter, body interface{}, code int) error {
	if body != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		w.WriteHeader(code)

		if err := json.NewEncoder(w).Encode(body); err != nil {
			return err
		}
	} else {
		w.WriteHeader(code)
	}

	return nil
}
