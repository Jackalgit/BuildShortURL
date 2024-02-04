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
		return nil, err
	}

	return &request, nil

}

func RequestListJSONToStruct(body io.Reader) ([]models.RequestBatch, error) {

	var requestList []models.RequestBatch

	dec := json.NewDecoder(body)
	if err := dec.Decode(&requestList); err != nil {
		return nil, err
	}

	return requestList, nil

}

func RequestListURLDelete(body io.Reader) ([]string, error) {

	var deleteList []string

	dec := json.NewDecoder(body)
	if err := dec.Decode(&deleteList); err != nil {
		return nil, err
	}

	return deleteList, nil

}
