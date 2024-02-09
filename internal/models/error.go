package models

type AddURLError struct {
	DupShortURLKey string
}

func NewAddURLError(dupShortURLKey string) error {
	return &AddURLError{DupShortURLKey: dupShortURLKey}
}

func (AD *AddURLError) Error() string {
	return AD.DupShortURLKey
}
