package main
import "encoding/json"

type peoplePayload struct {
	People      []string  `json:"people"`
	PayloadType   string  `json:"payloadType"`
}

type message struct {
	To          string 	 `json:"to"`
	From 		string 	 `json:from`
	Msg         string 	 `json:"msg"`
	PayloadType string   `json:"payloadType"`
}

func unmarshalMessage(b []byte) *message {
	var m *message = &message{}
	json.Unmarshal(b, m)
	return m
}