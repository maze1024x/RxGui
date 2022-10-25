package fatal

import (
    "os"
    "fmt"
    "strings"
)

func ThrowError(err error) {
    var desc = err.Error()
    fmt.Fprintf(os.Stderr, "error:\n%s\n", desc)
    if !(strings.HasSuffix(desc, "\n")) { fmt.Fprintf(os.Stderr, "\n") }
    os.Exit(126)
}
func ThrowBadArgsError(err error, help string) {
    fmt.Fprintf(os.Stderr, "bad command line arguments:\n%s\n\n", err.Error())
    fmt.Println(help)
    os.Exit(127)
}


