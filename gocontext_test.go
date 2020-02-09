// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0

package diag_test

import (
	"context"
	"testing"

	"github.com/urso/diag"
)

func TestStdlibContext(t *testing.T) {
	isContext := func(t *testing.T, ctx context.Context, dc *diag.Context) {
		t.Helper()
		has, ok := diag.DiagnosticsFrom(ctx)
		if !ok {
			t.Fatal("no diagnostic context found")
		}
		if dc != has {
			t.Fatal("unexpected diagnostic context returned")
		}
	}

	assertDiagnostics := func(t *testing.T, want map[string]interface{}, ctx context.Context) {
		t.Helper()
		dc, ok := diag.DiagnosticsFrom(ctx)
		if !ok {
			t.Fatal("no diagnostic context found")
		}
		assertCtx(t, want, dc)
	}

	t.Run("fail if no diagnostic context is stored", func(t *testing.T) {
		_, ok := diag.DiagnosticsFrom(context.Background())
		if ok {
			t.Fatal("context indicator should be false")
		}
	})

	t.Run("add diagnostics", func(t *testing.T) {
		dc := makeCtx(nil, nil, "a", 1)
		ctx := diag.NewDiagnostics(context.Background(), dc)
		isContext(t, ctx, dc)
	})

	t.Run("new context overwrites old", func(t *testing.T) {
		dc1 := makeCtx(nil, nil, "a", 1)
		dc2 := makeCtx(nil, nil, "a", 2)
		ctx := diag.NewDiagnostics(context.Background(), dc1)
		ctx = diag.NewDiagnostics(ctx, dc2)
		isContext(t, ctx, dc2)
	})

	t.Run("push creates a new context", func(t *testing.T) {
		ctx := diag.PushDiagnostics(context.Background(), "a", 1)
		assertDiagnostics(t, map[string]interface{}{"a": 1}, ctx)
	})

	t.Run("push references existing context", func(t *testing.T) {
		ctx := context.Background()
		ctx = diag.PushDiagnostics(ctx, "a", 1)
		ctx = diag.PushDiagnostics(ctx, "b", 2)
		assertDiagnostics(t, map[string]interface{}{"a": 1, "b": 2}, ctx)
	})

	t.Run("push shadows existing fields", func(t *testing.T) {
		ctx := context.Background()
		ctx = diag.PushDiagnostics(ctx, "a", 1)
		ctx = diag.PushDiagnostics(ctx, "a", 2)
		assertDiagnostics(t, map[string]interface{}{"a": 2}, ctx)
	})
}
