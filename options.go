package git

import (
	"fmt"
	"strings"
)

type Options []string

func (o Options) Want(hash string) string { return o.addheader(fmt.Sprintf("want %40s %s\n", hash, o)) }
func (o Options) Flush() string           { return o.addheader("") }
func (o Options) Done() string            { return o.addheader("done\n") }

func (o Options) addheader(v string) string {
	extra := 0
	if v != "" {
		extra = 4
	}
	return fmt.Sprintf("%04x%s", len(v)+extra, v)
}

func (o Options) String() string {
	return strings.Join(o, " ")
}
