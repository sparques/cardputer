package cardputer

import (
	"fmt"
	"io"
	"log"
)

var (
	KMsgWriter = io.Discard
	KMsgLogger *log.Logger
)

func kmsg(v ...any) {
	if KMsgLogger != nil {
		KMsgLogger.Print(v...)
		return
	}
	fmt.Fprint(KMsgWriter, v...)
}
