package cardputer

import (
	"fmt"
	"io"
	"log"
)

var (
	// KMsgWriter receives debug output when KMsgLogger is nil.
	KMsgWriter = io.Discard
	// KMsgLogger receives debug output when non-nil.
	KMsgLogger *log.Logger
)

func kmsg(v ...any) {
	if KMsgLogger != nil {
		KMsgLogger.Print(v...)
		return
	}
	fmt.Fprint(KMsgWriter, v...)
}
