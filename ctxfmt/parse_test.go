// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0

package ctxfmt

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type testHandler struct {
	t        *testing.T
	sequence []interface{}
}

func TestParser(t *testing.T) {
	cases := map[string][]interface{}{
		"simple string": {
			"simple string",
		},
		"double %% in middle": {
			"double %",
			" in middle",
		},
		"double %%": {
			"double %",
		},
		"%v": {
			formatToken{verb: 'v'},
		},
		"verb %v": {
			"verb ",
			formatToken{verb: 'v'},
		},
		"verb %v in middle": {
			"verb ",
			formatToken{verb: 'v'},
			" in middle",
		},
		"%+v": {
			formatToken{verb: 'v', flags: flags{plusV: true}},
		},
		"%#v": {
			formatToken{verb: 'v', flags: flags{sharpV: true}},
		},
		"%-v": {
			formatToken{verb: 'v', flags: flags{minus: true, zero: false}},
		},
		"%0v": {
			formatToken{verb: 'v', flags: flags{zero: true}},
		},
		"%-0v": {
			formatToken{verb: 'v', flags: flags{minus: true, zero: false}},
		},
		"% v": {
			formatToken{verb: 'v', flags: flags{space: true}},
		},
		"%5d": {
			formatToken{verb: 'd', width: 5, flags: flags{hasWidth: true}},
		},
		"%.3f": {
			formatToken{verb: 'f', precision: 3, flags: flags{hasPrecision: true}},
		},
		"%5.3f": {
			formatToken{verb: 'f', width: 5, precision: 3, flags: flags{hasWidth: true, hasPrecision: true}},
		},
		"%.d": {
			formatToken{verb: 'd', flags: flags{hasPrecision: true}},
		},
		"unknown verb %a": {
			"unknown verb ",
			errInvalidVerb,
		},
		"no verb %": {
			"no verb ",
			errNoVerb,
		},
		"%12": {
			errNoVerb,
		},
		"%{field}": {
			formatToken{verb: 'v', field: "field", flags: flags{named: true}},
		},
		"%{": {
			errCloseMissing,
		},
		"%{oops": {
			errCloseMissing,
		},
		"%{oops:v": {
			errCloseMissing,
		},
		"at end %{field}": {
			"at end ",
			formatToken{verb: 'v', field: "field", flags: flags{named: true}},
		},
		"%{field} at the beginning": {
			formatToken{verb: 'v', field: "field", flags: flags{named: true}},
			" at the beginning",
		},
		"field %{name} in the middle": {
			"field ",
			formatToken{verb: 'v', field: "name", flags: flags{named: true}},
			" in the middle",
		},
		"%{field:d}": {
			formatToken{verb: 'd', field: "field", flags: flags{named: true}},
		},
		"%{field:+d}": {
			formatToken{verb: 'd', field: "field", flags: flags{plus: true, named: true}},
		},
		"%{field:#d}": {
			formatToken{verb: 'd', field: "field", flags: flags{sharp: true, named: true}},
		},
		"%{field:5d}": {
			formatToken{verb: 'd', field: "field", width: 5, flags: flags{hasWidth: true, named: true}},
		},
		"%{field:.3d}": {
			formatToken{verb: 'd', field: "field", precision: 3, flags: flags{hasPrecision: true, named: true}},
		},
		"%{field:5.3d}": {
			formatToken{verb: 'd', field: "field", width: 5, precision: 3, flags: flags{hasWidth: true, hasPrecision: true, named: true}},
		},
		"%{+field}": {
			formatToken{verb: 'v', field: "field", flags: flags{plusV: true, named: true}},
		},
		"%{#field}": {
			formatToken{verb: 'v', field: "field", flags: flags{sharpV: true, named: true}},
		},
		"%{field:a}": {
			errInvalidVerb,
		},
	}

	for str, want := range cases {
		t.Run(str, func(t *testing.T) {
			handler := &testHandler{t: t}
			p := parser{handler: handler}
			p.parse(str)

			opts := []cmp.Option{
				cmp.Comparer(func(a, b formatToken) bool {
					return a == b
				}),
				cmp.Comparer(func(a, b error) bool {
					return a == b
				}),
			}

			if diff := cmp.Diff(want, handler.sequence, opts...); diff != "" {
				t.Fatalf("missmatch (-want +got):\n%s", diff)
			}
		})
	}
}

func (t *testHandler) onString(s string)                     { t.sequence = append(t.sequence, s) }
func (t *testHandler) onToken(tok formatToken)               { t.sequence = append(t.sequence, tok) }
func (t *testHandler) onParseError(_ formatToken, err error) { t.sequence = append(t.sequence, err) }
