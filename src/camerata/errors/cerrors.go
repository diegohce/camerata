package errors

type CamerataError struct {
	Message string
}

func (me CamerataError) Error() string {
	return me.Message
}

type CamerataArgumentsError CamerataError

func (me CamerataArgumentsError) Error() string {
	return me.Message
}
