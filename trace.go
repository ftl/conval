package conval

import "log"

var traceActive bool = false

func tracef(format string, args ...any) {
	if !traceActive {
		return
	}
	log.Printf(format, args...)
}
