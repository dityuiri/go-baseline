package controller

// import (
// 	"encoding/json"
//
// 	"github.com/dityuiri/go-baseline/model"
// 	"github.com/dityuiri/go-baseline/service"
// )
//
// type ConsumerHandler struct {
// 	TransactionFeed service.ITransactionFeedService
// }
//
// func (ch *ConsumerHandler) Transaction(msg []byte) (bool, error) {
// 	var transaction = &model.Transaction{}
//
// 	err := json.Unmarshal(msg, transaction)
// 	if err != nil {
// 		return true, err
// 	}
//
// 	return ch.TransactionFeed.TransactionRecorded(transaction)
// }
