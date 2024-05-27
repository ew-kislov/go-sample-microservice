package api

type Error struct {
	Code    int
	Message string
}

func (err Error) Error() string {
	return err.Message
}
