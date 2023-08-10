package rescue

import (
	"fmt"
	"runtime/debug"

	"github.com/mcdull-kk/pkg/log"
)

// Recover is used with defer to do cleanup on panics.
// Use it like:
//  defer Recover(func() {})
func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		log.Error("panic", fmt.Sprintf("%v\n%s", p, string(debug.Stack())))
	}
}

func GoSafe(fn func()) {
	go func() {
		defer Recover()
		fn()
	}()
}
