package infrastructure

type ApiError struct {
	Code    int16
	Message string
	Details string
}

func (err ApiError) Error() string {
	return err.Message
}
