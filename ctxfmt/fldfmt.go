package ctxfmt

import (
	"io"
	"strings"
)

type CB func(key string, idx int, val interface{})

func Sprintf(cb CB, msg string, vs ...interface{}) (string, []interface{}) {
	var buf strings.Builder
	rest := Fprintf(&buf, cb, msg, vs...)
	return buf.String(), rest
}

func Fprintf(w io.Writer, cb CB, msg string, vs ...interface{}) (rest []interface{}) {
	in := &interpreter{
		cb:   cb,
		p:    &printer{To: w},
		args: argstate{args: vs},
	}
	parser := &parser{handler: in}
	parser.parse(msg)

	used := in.args.idx
	if used >= len(vs) {
		return nil
	}

	// collect errors from extra variables
	rest = vs[used:]
	for i := range rest {
		if isErrorValue(rest[i]) {
			cb("", used+i, rest[i])
		}
	}
	return rest
}
