package assert

import "fmt"

func Assert(condition bool, str ...string) {
	if condition {
		return
	} else {
		panicMsg := ""
		if len(str) > 0 {
			panicMsg = fmt.Sprintln(str)
		} else {
			panicMsg = "Assertion failed\n"
		}

		panic(panicMsg)
	}
}
