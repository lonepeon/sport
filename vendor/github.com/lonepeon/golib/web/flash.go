package web

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(FlashMessage{})
}

type FlashMessage struct {
	Kind    string
	Message string
}

func NewFlashMessageError(pattern string, vars ...interface{}) FlashMessage {
	return FlashMessage{Kind: "error", Message: fmt.Sprintf(pattern, vars...)}
}

func NewFlashMessageSuccess(pattern string, vars ...interface{}) FlashMessage {
	return FlashMessage{Kind: "success", Message: fmt.Sprintf(pattern, vars...)}
}
