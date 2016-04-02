package main

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

type CamerataConnectionError CamerataError

func (me CamerataConnectionError) Error() string {
	return me.Message
}

type CamerataRunError CamerataError

func (me CamerataRunError) Error() string {
	return me.Message
}

type CamerataModuleError CamerataError

func (me CamerataModuleError) Error() string {
	return me.Message
}
