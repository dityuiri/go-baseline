package common

import "encoding/json"

// Alias for function patch that can't be easily mocked
var (
	JsonMarshal   = json.Marshal
	JsonUnmarshal = json.Unmarshal
)
