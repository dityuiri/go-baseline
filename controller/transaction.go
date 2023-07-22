package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"stockbit-challenge/service"
)

const (
	maxMultipartForm          = 32 << 20
	defaultMultiPartFormQuery = "file"
	fileNdJSON                = "ndjson"
)

type (
	ITransactionController interface {
		UploadTransactions(w http.ResponseWriter, r *http.Request)
	}

	TransactionController struct {
		TransactionFeedService service.ITransactionFeedService
	}
)

func (c *TransactionController) UploadTransactions(w http.ResponseWriter, r *http.Request) {
	buff, err := c.retrieveFile(r, defaultMultiPartFormQuery, fileNdJSON)
	if err != nil {
		c.writeResponse(w, http.StatusBadRequest, err)
		return
	}

	err = c.TransactionFeedService.ProduceTransaction(buff)
	if err != nil {
		c.writeResponse(w, http.StatusInternalServerError, err)
		return
	}

	c.writeResponse(w, http.StatusNoContent, "")
}

func (*TransactionController) retrieveFile(r *http.Request, query string, extFile string) (bytes.Buffer, error) {
	var buff bytes.Buffer

	if query == "" {
		query = defaultMultiPartFormQuery
	}

	if err := r.ParseMultipartForm(maxMultipartForm); err != nil {
		return buff, err
	}

	file, handler, err := r.FormFile(query)
	if err != nil {
		return buff, err
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	var (
		filenameExt   string
		splitFilename = strings.Split(handler.Filename, ".")
	)

	if len(splitFilename) > 0 {
		filenameExt = splitFilename[len(splitFilename)-1]
	}

	if filenameExt != extFile {
		return buff, fmt.Errorf("expecting '%s', got '%s'", extFile, filenameExt)
	}

	if _, err = io.Copy(&buff, file); err != nil {
		return buff, err
	}

	return buff, nil
}

func (*TransactionController) writeResponse(w http.ResponseWriter, httpStatus int, response interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)

	resp, _ := json.Marshal(response)
	_, _ = w.Write(resp)
}
