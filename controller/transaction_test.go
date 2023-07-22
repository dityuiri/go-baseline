package controller

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	serviceMock "stockbit-challenge/mock/service"
)

var _ = Describe("TransactionController", func() {
	var (
		path = "/publish/transaction"

		mockCtrl                   *gomock.Controller
		mockTransactionFeedService *serviceMock.MockITransactionFeedService

		controller TransactionController
		router     *chi.Mux
		server     *ghttp.Server
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockTransactionFeedService = serviceMock.NewMockITransactionFeedService(mockCtrl)

		controller = TransactionController{
			TransactionFeedService: mockTransactionFeedService,
		}

		router = chi.NewRouter()
		router.Post(path, controller.UploadTransactions)

		server = ghttp.NewServer()
	})

	AfterEach(func() {
		mockCtrl.Finish()
		server.Close()
	})

	Context("Success", func() {
		It("should return a successful response", func() {
			file, _ := os.Open("test.ndjson")
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.ndjson")
			_, _ = io.Copy(part, file)
			_ = writer.Close()

			request, _ := http.NewRequest(http.MethodPost, server.URL()+path, body)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			mockTransactionFeedService.EXPECT().ProduceTransaction(gomock.Any()).Return(nil)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, request)

			Expect(resp.Code).To(Equal(http.StatusNoContent))
		})
	})

	Context("Internal server error from service", func() {
		It("should return an internal server error response", func() {
			file, _ := os.Open("test.ndjson")
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.ndjson")
			_, _ = io.Copy(part, file)
			_ = writer.Close()

			request, _ := http.NewRequest(http.MethodPost, server.URL()+path, body)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			mockTransactionFeedService.EXPECT().ProduceTransaction(gomock.Any()).Return(errors.New("error"))

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, request)

			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("Error retrieving file", func() {
		It("should return a bad request response", func() {
			file, _ := os.Open("test.perkedel")
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", "test.perkedel")
			_, _ = io.Copy(part, file)
			_ = writer.Close()

			request, _ := http.NewRequest(http.MethodPost, server.URL()+path, body)
			request.Header.Set("Content-Type", writer.FormDataContentType())

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, request)

			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
