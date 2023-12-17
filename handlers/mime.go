package handlers

import (
	"fmt"
	"net/http"

	"github.com/mholzen/information/triples"
)

type MimeType string

func (m MimeType) String() string {
	return string(m)
}

func (m MimeType) LessThan(other triples.Node) bool {
	return triples.NodeLessThan(m, other)
}

func NewMimeType(input Payload) (MimeType, error) {
	switch data := input.Data.(type) {
	case FileInfo:
		res, err := data.ContentType()
		if err != nil {
			return "", err
		}
		return MimeType(res), nil
	case string:
		return MimeType(http.DetectContentType([]byte(data))), nil
	default:
		return "", fmt.Errorf("cannot convert '%T' to MimeType", input.Data)
	}
}

func ToMimeType(input Payload) (Payload, error) {
	mimeType, err := NewMimeType(input)
	if err != nil {
		return input, err
	}

	return Payload{
		Content: "text/plain",
		Data:    mimeType,
	}, nil
}
