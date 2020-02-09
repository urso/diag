// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0

package ctxfmt

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSprintfFields(t *testing.T) {
	type cbRecord struct {
		Key string
		Idx int
		Val interface{}
	}
	type records []cbRecord

	values := func(vs ...interface{}) []interface{} { return vs }

	cases := []struct {
		in   string
		out  string
		args []interface{}
		want records
		rest []interface{}
	}{
		{
			in:  "hello world",
			out: "hello world",
		},
		{
			in:   "%{field}",
			out:  "3",
			args: values(3),
			want: records{
				{"field", 0, 3},
			},
		},
		{
			in:   "%{a}%{b}%{c}",
			out:  "123",
			args: values(1, 2, 3),
			want: records{
				{"a", 0, 1},
				{"b", 1, 2},
				{"c", 2, 3},
			},
		},
		{
			in:   "%{a}%03v%{c}",
			out:  "10023",
			args: values(1, 2, 3),
			want: records{
				{"a", 0, 1},
				{"c", 2, 3},
			},
		},
		{
			in:   "%{field} %{noarg}",
			out:  "3 %!v(MISSING)",
			args: values(3),
			want: records{
				{"field", 0, 3},
			},
		},
		{
			in:   "%{field}",
			out:  "1",
			args: values(1, 2, 3),
			want: records{
				{"field", 0, 1},
			},
			rest: values(2, 3),
		},
	}

	for i, test := range cases {
		name := fmt.Sprintf("%d: %v -> %v", i, test.in, test.out)
		t.Run(name, func(t *testing.T) {
			var actual records
			out, rest := Sprintf(func(key string, idx int, val interface{}) {
				actual = append(actual, cbRecord{key, idx, val})
			}, test.in, test.args...)

			if test.out != out {
				t.Errorf("Output failure. Want <%s>, Got <%s>", test.out, out)
			}

			if diff := cmp.Diff(test.want, actual); diff != "" {
				t.Errorf("callback missmatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(test.rest, rest); diff != "" {
				t.Errorf("rest fields missmatch (-want +got):\n%s", diff)
			}
		})
	}
}
