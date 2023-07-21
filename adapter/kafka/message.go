package kafka

type Message struct {
	Partition int
	Offset    int64

	Key   []byte
	Value interface{}

	Headers Header
}

type Header map[string][]byte
