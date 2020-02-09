// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0

package diag_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/urso/diag"
)

type testVisitor struct {
	M     map[string]interface{}
	keys  []string
	stack []map[string]interface{}
}

func makeCtx(before, after *diag.Context, vs ...interface{}) *diag.Context {
	ctx := diag.NewContext(before, after)
	ctx.AddAll(vs...)
	return ctx
}

func assertCtx(t *testing.T, want map[string]interface{}, ctx *diag.Context) {
	t.Helper()

	var v testVisitor
	ctx.VisitStructured(&v)
	requireEqual(t, want, v.Get())
}

func assertFlatCtx(t *testing.T, want map[string]interface{}, ctx *diag.Context) {
	t.Helper()

	var v testVisitor
	ctx.VisitKeyValues(&v)
	requireEqual(t, want, v.Get())
}

func requireEqual(t *testing.T, want, has interface{}, msg ...string) {
	t.Helper()
	if diff := cmp.Diff(want, has); diff != "" {
		if len(msg) > 0 {
			t.Fatalf("%v -- missmatch (-want +got):\n%s", diff, strings.Join(msg, " "))
		}
		t.Fatalf("missmatch (-want +got):\n%s", diff)
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("require no error fail with: %+v", err)
	}
}

func (v *testVisitor) Get() map[string]interface{} {
	return v.M
}

func (v *testVisitor) OnValue(key string, val diag.Value) error {
	if v.M == nil {
		v.M = map[string]interface{}{}
	}

	v.M[key] = val.Interface()
	return nil
}

func (v *testVisitor) OnObjStart(key string) error {
	v.keys = append(v.keys, key)
	v.stack = append(v.stack, v.M)
	v.M = nil
	return nil
}

func (v *testVisitor) OnObjEnd() error {
	keysEnd := len(v.keys) - 1
	key := v.keys[keysEnd]
	v.keys = v.keys[:keysEnd]

	m := v.M
	mapsEnd := len(v.stack) - 1
	v.M = v.stack[mapsEnd]
	v.stack = v.stack[:mapsEnd]

	if m != nil {
		if v.M == nil {
			v.M = map[string]interface{}{}
		}
		v.M[key] = m
	}

	return nil
}
