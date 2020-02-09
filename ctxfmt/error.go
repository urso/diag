package ctxfmt

import "errors"

var (
	errInvalidVerb  = errors.New("invalid verb")
	errNoVerb       = errors.New("no verb")
	errCloseMissing = errors.New("missing '}'")
	errNoFieldName  = errors.New("field name missing")
	errMissingArg   = errors.New("missing arg")
)
