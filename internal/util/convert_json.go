package util

import (
	"bytes"
	"encoding/json"
	"github.com/Jackalgit/BuildShortURL/internal/models"
	"io"
)

func RequestJSONToStruct(body io.Reader) (*models.Request, error) {
	var request models.Request

	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)
	if err != nil {
		return &request, err
	}
	// десериализуем JSON в Request
	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		return &request, err
	}

	return &request, nil

}
