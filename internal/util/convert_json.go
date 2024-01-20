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

func RequestListJSONToStruct(body io.Reader) ([]models.RequestBatch, error) {

	var requestList []models.RequestBatch

	dec := json.NewDecoder(body)
	if err := dec.Decode(&requestList); err != nil {
		return requestList, err
	}

	//var buf bytes.Buffer
	//_, err := buf.ReadFrom(body)
	//if err != nil {
	//	return requestList, err
	//}
	//// десериализуем JSON в RequestList
	//if err = json.Unmarshal(buf.Bytes(), &requestList); err != nil {
	//	return requestList, err
	//}

	return requestList, nil

}
