package api

type ApiError struct {
	Code    int
	Message string
}

func (err ApiError) Error() string {
	return err.Message
}
