package model

type ConnectCmspError struct {
}

func (ConnectCmspError) Error() string {
	return "Connect cmsp Error"
}
