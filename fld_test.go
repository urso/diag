// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0

package diag_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/urso/diag"
)

func TestFieldValues(t *testing.T) {
	now := time.Now()

	cases := map[diag.Field]interface{}{
		diag.Bool("b", true):                             true,
		diag.Bool("b", false):                            false,
		diag.Int("i", -23):                               -23,
		diag.Int64("i", -42):                             int64(-42),
		diag.Uint("i", 23):                               uint64(23),
		diag.Uint64("i", math.MaxInt32):                  uint64(math.MaxInt32),
		diag.Float("f", 3.14):                            3.14,
		diag.String("hello", "world"):                    "world",
		diag.Duration("d", time.Duration(5*time.Second)): 5 * time.Second,
		diag.Timestamp("ts", now):                        now,
	}

	var i int
	for field, want := range cases {
		name := fmt.Sprintf("%v: %#v", i, field)
		t.Run(name, func(t *testing.T) {
			if diff := cmp.Diff(want, field.Value.Interface()); diff != "" {
				t.Errorf("missmatch (+want, -got):\n%v", diff)
			}
		})
	}
}
