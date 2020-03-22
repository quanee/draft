package draft

import "fmt"

func debugPrint(format string, values ...interface{}) {
	fmt.Fprintf(DefaultWriter, format, values...)
}
