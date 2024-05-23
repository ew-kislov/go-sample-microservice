package api

type ApiError struct {
	Code    int16
	Message string
}

func (err ApiError) Error() string {
	return err.Message
}
