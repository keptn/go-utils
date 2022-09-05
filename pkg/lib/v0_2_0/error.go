package v0_2_0

type Error struct {
	StatusType StatusType
	ResultType ResultType
	Message    string
	Err        error
}

func (e Error) Error() string {
	return e.Message
}
