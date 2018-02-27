package model

type NoNeedQueryError struct {
}

func (NoNeedQueryError) Error() string {
	return "No Need Query"
}
