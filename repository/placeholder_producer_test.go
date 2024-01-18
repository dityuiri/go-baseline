package repository

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	producerMock "github.com/dityuiri/go-baseline/adapter/kafka/producer/mock"
	"github.com/dityuiri/go-baseline/config"
	"github.com/dityuiri/go-baseline/model"
)

func TestPlaceholderProducer_ProducePlaceholderRecord(t *testing.T) {
	var (
		mockCtrl     = gomock.NewController(t)
		mockProducer = producerMock.NewMockIProducer(mockCtrl)

		producer = PlaceholderProducer{
			Producer: mockProducer,
			KafkaConfig: &config.Kafka{
				ProducerTopics: map[string]string{
					"placeholder":     "placeholder",
					"placeholder_dlq": "placeholder_dlq",
				},
			},
		}

		ctx                = context.Background()
		placeholderMessage = model.PlaceholderMessage{
			ID: uuid.New().String(),
		}
	)

	t.Run("positive", func(t *testing.T) {
		mockProducer.EXPECT().Produce(ctx, "placeholder", gomock.Any())

		err := producer.ProducePlaceholderRecord(ctx, placeholderMessage)
		assert.Nil(t, err)
	})

}
