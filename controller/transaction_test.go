package controller

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"

	"stockbit-challenge/mock"
	serviceMock "stockbit-challenge/mock/service"
)

func TestTransactionController_UploadTransactions(t *testing.T) {
	var (
		path = "/publish/transaction"

		mockCtrl                   = gomock.NewController(t)
		mockResp                   = mock.NewMockResponseWriter(mockCtrl)
		mockTransactionFeedService = serviceMock.NewMockITransactionFeedService(mockCtrl)

		controller = TransactionController{
			TransactionFeedService: mockTransactionFeedService,
		}
	)

	router := chi.NewRouter()
	router.Post(path, controller.UploadTransactions)
	defer mockCtrl.Finish()

	t.Run("success", func(t *testing.T) {
		file, _ := os.Open("test.ndjson")
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.ndjson")
		_, _ = io.Copy(part, file)
		_ = writer.Close()

		request, _ := http.NewRequest(http.MethodPost, path, body)
		request.Header.Set("Content-Type", writer.FormDataContentType())

		mockTransactionFeedService.EXPECT().ProduceTransaction(gomock.Any()).Return(nil)

		mockResp.EXPECT().Header().Return(http.Header{})
		mockResp.EXPECT().WriteHeader(http.StatusNoContent)
		mockResp.EXPECT().Write(gomock.Any()).Return(204, nil)

		router.ServeHTTP(mockResp, request)
	})

	t.Run("internal server error from service", func(t *testing.T) {
		file, _ := os.Open("test.ndjson")
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.ndjson")
		_, _ = io.Copy(part, file)
		_ = writer.Close()

		request, _ := http.NewRequest(http.MethodPost, path, body)
		request.Header.Set("Content-Type", writer.FormDataContentType())

		mockTransactionFeedService.EXPECT().ProduceTransaction(gomock.Any()).Return(errors.New("error"))

		mockResp.EXPECT().Header().Return(http.Header{})
		mockResp.EXPECT().WriteHeader(http.StatusInternalServerError)
		mockResp.EXPECT().Write(gomock.Any()).Return(500, nil)

		router.ServeHTTP(mockResp, request)
	})

	t.Run("error retrieving file", func(t *testing.T) {
		file, _ := os.Open("test.perkedel")
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test.perkedel")
		_, _ = io.Copy(part, file)
		_ = writer.Close()

		request, _ := http.NewRequest(http.MethodPost, path, body)
		request.Header.Set("Content-Type", writer.FormDataContentType())

		mockResp.EXPECT().Header().Return(http.Header{})
		mockResp.EXPECT().WriteHeader(http.StatusBadRequest)
		mockResp.EXPECT().Write(gomock.Any()).Return(400, nil)

		router.ServeHTTP(mockResp, request)
	})
}
