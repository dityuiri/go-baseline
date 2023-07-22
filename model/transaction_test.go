package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RawTransaction to Transaction", func() {
	var (
		rawTypeE RawTransaction
	)

	BeforeEach(func() {
		rawTypeE = RawTransaction{
			Type:             "E",
			OrderNumber:      "12831289312389",
			OrderVerb:        "S",
			Quantity:         "13",
			Price:            "20000",
			StockCode:        "BBCA",
			ExecutionPrice:   "20000",
			ExecutedQuantity: "13",
			OrderBook:        "123",
		}
	})

	Context("Positive", func() {
		It("should convert rawTypeE to transaction successfully", func() {
			trx, err := rawTypeE.ToTransaction()
			Expect(trx).ToNot(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Invalid quantity", func() {
		It("should return an error", func() {
			rawTypeInvalid := RawTransaction{
				Type:             "E",
				OrderNumber:      "12831289312389",
				OrderVerb:        "S",
				Quantity:         "13",
				Price:            "20000",
				StockCode:        "BBCA",
				ExecutionPrice:   "20000",
				ExecutedQuantity: "HAAAA",
				OrderBook:        "123",
			}

			trx, err := rawTypeInvalid.ToTransaction()
			Expect(trx.StockCode).To(BeIdenticalTo("BBCA"))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Invalid price", func() {
		It("should return an error", func() {
			rawTypeInvalid := RawTransaction{
				Type:             "E",
				OrderNumber:      "12831289312389",
				OrderVerb:        "S",
				Quantity:         "13",
				Price:            "20000",
				StockCode:        "BBCA",
				ExecutionPrice:   "Sachi",
				ExecutedQuantity: "13",
				OrderBook:        "123",
			}

			trx, err := rawTypeInvalid.ToTransaction()
			Expect(trx.StockCode).To(BeIdenticalTo("BBCA"))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Invalid order book", func() {
		It("should return an error", func() {
			rawTypeInvalid := RawTransaction{
				Type:             "E",
				OrderNumber:      "12831289312389",
				OrderVerb:        "S",
				Quantity:         "13",
				Price:            "20000",
				StockCode:        "BBCA",
				ExecutionPrice:   "20000",
				ExecutedQuantity: "13",
				OrderBook:        "JASDASD",
			}

			trx, err := rawTypeInvalid.ToTransaction()
			Expect(trx.StockCode).To(BeEmpty())
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Transaction IsAccountable", func() {
	Context("True", func() {
		It("should return true", func() {
			var trx = Transaction{Type: "E"}
			res := trx.IsAccountable()
			Expect(res).To(BeTrue())
		})
	})

	Context("False", func() {
		It("should return false", func() {
			var trx = Transaction{Type: "BLAU"}
			res := trx.IsAccountable()
			Expect(res).To(BeFalse())
		})
	})
})

var _ = Describe("Transaction IsPreviousPrice", func() {
	Context("True", func() {
		It("should return true", func() {
			var trx = Transaction{Type: "E"}
			res := trx.IsPreviousPrice()
			Expect(res).To(BeTrue())
		})
	})

	Context("False", func() {
		It("should return false", func() {
			var trx = Transaction{Type: "BLAU"}
			res := trx.IsAccountable()
			Expect(res).To(BeFalse())
		})
	})
})
