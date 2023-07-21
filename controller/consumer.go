package controller

import (
	"encoding/json"

	"stockbit-challenge/model"
	"stockbit-challenge/service"
)

type ConsumerHandler struct {
	TransactionFeed service.ITransactionFeedService
}

func (ch *ConsumerHandler) Transaction(msg []byte) (bool, error) {
	var transaction = &model.Transaction{}

	err := json.Unmarshal(msg, transaction)
	if err != nil {
		return true, err
	}

	return ch.TransactionFeed.TransactionRecorded(transaction)
}
