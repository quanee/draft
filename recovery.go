package draft

import (
	//"draft/utils/log"
	"fmt"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
)

// print stack trace for debug
func stack(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				httpRequest, _ := httputil.DumpRequest(c.Req, false)
				//log.Printf("%s\n%s", stack(message), string(httpRequest))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
