package errors

import (
	"runtime"
	"fmt"
)

const AccessDenied = "Access Denied"

type Error struct {
	Message string
	Pc []uintptr
}

func (this Error) Error() string {
	output := this.Message + "\n"
	frames := runtime.CallersFrames(this.Pc)
	frame, more := frames.Next()
	for more {
		output = fmt.Sprint(output, "\n"," ", frame.Function,"\n\t ", frame.File,":", frame.Line)
		if (more) {
			frame, more = frames.Next()
		}
	}
	return output
}

func New(input string) error {
	var output Error
	output.Message = input
	output.Pc = make([]uintptr, 20)
	runtime.Callers(2, output.Pc)
	return output
}