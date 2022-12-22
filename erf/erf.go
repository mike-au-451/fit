package erf

import (
	"fmt"
	"runtime"
)

// func Fatal(format string, args ...interface{}) {
// 	_, fn, ln, ok := runtime.Caller(1)
// 	if !ok {
// 		fn = "unknown"
// 		ln = 0
// 	}
// 	where := fmt.Sprintf("FATAL: %s:%d", fn, ln)
// 	msg := fmt.Sprintf(format, args)

// 	fmt.Printf("%s: %s\n", where, msg)
// }

func Here() (string, int) {
	_, fn, ln, ok := runtime.Caller(1)
	if !ok {
		fn = "(stack corrupt)"
		ln = 0
	}

	return fn, ln
}

func Errorf(format string, args ...interface{}) error {
	_, fn, ln, ok := runtime.Caller(1)
	if !ok {
		fn = "(stack corrupt)"
		ln = 0
	}

 	where := fmt.Sprintf("FATAL: %s:%d", fn, ln)
	msg := fmt.Sprintf(format, args...)

	return fmt.Errorf("%s: %s\n", where, msg)
}
