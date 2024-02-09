package models

type StatusURL struct {
	Found  bool
	Delete bool
}

func NewStatusURL(found bool, delete bool) StatusURL {
	return StatusURL{Found: found, Delete: delete}
}
