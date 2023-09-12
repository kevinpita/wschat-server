package message

import "github.com/brianvoe/gofakeit/v6"

type Message struct {
	Data int `json:"data"`
}

func Random() (*Message, error) {
	var msg Message
	errFake := gofakeit.Struct(&msg)
	if errFake != nil {
		return nil, errFake
	}
	return &msg, nil
}
