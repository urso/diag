package diag_test

import (
	"testing"

	"github.com/urso/diag"
)

func TestCtxBuild(t *testing.T) {
	t.Run("new empty context", func(t *testing.T) {

		ctx := diag.NewContext(nil, nil)
		assertCtx(t, nil, ctx)
	})

	t.Run("new empty with non-empty before", func(t *testing.T) {
		before := diag.NewContext(nil, nil)
		before.AddAll("hello", "world")
		ctx := diag.NewContext(before, nil)
		assertCtx(t, map[string]interface{}{
			"hello": "world",
		}, ctx)
	})

	t.Run("new empty with non-empty after", func(t *testing.T) {
		after := diag.NewContext(nil, nil)
		after.AddAll("hello", "world")
		ctx := diag.NewContext(nil, after)
		assertCtx(t, map[string]interface{}{
			"hello": "world",
		}, ctx)
	})

	t.Run("new empty with non-empty before and after", func(t *testing.T) {
		before := diag.NewContext(nil, nil)
		before.AddAll("before", "hello", "overwrite", 1)

		after := diag.NewContext(nil, nil)
		after.AddAll("after", "world", "overwrite", 2)

		ctx := diag.NewContext(before, after)
		assertCtx(t, map[string]interface{}{
			"before":    "hello",
			"after":     "world",
			"overwrite": 2,
		}, ctx)
	})

	t.Run("new context overwrites before elements", func(t *testing.T) {
		before := diag.NewContext(nil, nil)
		before.AddAll("before", "hello", "overwrite", 1)

		ctx := diag.NewContext(before, nil)
		ctx.AddAll("overwrite", 2)
		assertCtx(t, map[string]interface{}{
			"before":    "hello",
			"overwrite": 2,
		}, ctx)
	})

	t.Run("new context does not overwrite before elements", func(t *testing.T) {
		after := diag.NewContext(nil, nil)
		after.AddAll("hello", "world", "overwrite", 1)

		ctx := diag.NewContext(nil, after)
		ctx.AddAll("overwrite", 2)
		assertCtx(t, map[string]interface{}{
			"hello":     "world",
			"overwrite": 1,
		}, ctx)
	})
}

func TestCtxAdd(t *testing.T) {
	ctx := diag.NewContext(nil, nil)
	ctx.Add("hello", diag.ValString("world"))
	assertCtx(t, map[string]interface{}{
		"hello": "world",
	}, ctx)
}

func TestCtxAddAll(t *testing.T) {
	cases := map[string]struct {
		in   []interface{}
		want map[string]interface{}
	}{
		"unique keys": {
			in:   []interface{}{"key1", 1, "key2", 2},
			want: map[string]interface{}{"key1": 1, "key2": 2},
		},
		"duplicate keys": {
			in:   []interface{}{"key", 1, "key", 2},
			want: map[string]interface{}{"key": 2},
		},
		"accepts Value": {
			in:   []interface{}{"key", diag.ValInt(10)},
			want: map[string]interface{}{"key": 10},
		},
		"accepts Field": {
			in:   []interface{}{diag.Field{Key: "key", Value: diag.ValInt(10)}},
			want: map[string]interface{}{"key": 10},
		},
		"mix fields with key values": {
			in: []interface{}{
				"before", "hello",
				diag.Field{Key: "key", Value: diag.ValInt(2)},
				"after", "world",
			},
			want: map[string]interface{}{
				"before": "hello",
				"key":    2,
				"after":  "world",
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := diag.NewContext(nil, nil)
			ctx.AddAll(test.in...)
			assertCtx(t, test.want, ctx)
		})
	}
}

func TestCtxAddField(t *testing.T) {
	t.Run("user field", func(t *testing.T) {
		ctx := diag.NewContext(nil, nil)
		ctx.AddField(diag.String("hello", "world"))
		assertCtx(t, map[string]interface{}{
			"hello": "world",
		}, ctx)
		requireEqual(t, 1, ctx.User().Len())
		requireEqual(t, 0, ctx.Standardized().Len())
	})

	t.Run("standardized field", func(t *testing.T) {
		f := diag.String("hello", "world")
		f.Standardized = true
		ctx := diag.NewContext(nil, nil)
		ctx.AddField(f)
		assertCtx(t, map[string]interface{}{
			"hello": "world",
		}, ctx)
		requireEqual(t, 0, ctx.User().Len())
		requireEqual(t, 1, ctx.Standardized().Len())
	})
}

func TestCtxAddFields(t *testing.T) {
	cases := map[string]struct {
		in        []diag.Field
		want      map[string]interface{}
		user, std int
	}{
		"unique keys": {
			in: []diag.Field{
				diag.Int("key1", 1),
				diag.Int("key2", 2),
			},
			want: map[string]interface{}{"key1": 1, "key2": 2},
			user: 2,
		},
		"duplicate keys": {
			in: []diag.Field{
				diag.Int("key", 1),
				diag.Int("key", 2),
			},
			want: map[string]interface{}{"key": 2},
			user: 2, // both keys are stored
		},
		"standardized and user fields": {
			in: []diag.Field{
				diag.Int("key", 1),
				diag.Field{Key: "test", Value: diag.ValInt(2), Standardized: true},
			},
			want: map[string]interface{}{"key": 1, "test": 2},
			user: 1,
			std:  1,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := diag.NewContext(nil, nil)
			ctx.AddFields(test.in...)
			assertCtx(t, test.want, ctx)
			requireEqual(t, test.user, ctx.User().Len())
			requireEqual(t, test.std, ctx.Standardized().Len())
		})
	}
}

func TestCtxFiltered(t *testing.T) {
	filtLocal := (*diag.Context).Local
	filtUser := (*diag.Context).User
	filtStd := (*diag.Context).Standardized

	cases := map[string]struct {
		in     *diag.Context
		filter func(*diag.Context) *diag.Context
		want   map[string]interface{}
		len    int // number of entries in new context
	}{
		"local filter ignores before": {
			in: makeCtx(makeCtx(nil, nil, "before", "hello"), nil,
				"current", "world"),
			filter: filtLocal,
			want: map[string]interface{}{
				"current": "world",
			},
			len: 1,
		},
		"local filter ignores after": {
			in: makeCtx(nil, makeCtx(nil, nil, "after", "world"),
				"key", "value"),
			filter: filtLocal,
			want: map[string]interface{}{
				"key": "value",
			},
			len: 1,
		},

		"user filter transitive": {
			in: makeCtx(
				makeCtx(nil, nil, diag.String("user_before", "test"),
					diag.Field{Key: "std_before", Value: diag.ValInt(1), Standardized: true}),
				makeCtx(nil, nil, diag.String("user_after", "test"),
					diag.Field{Key: "std_after", Value: diag.ValInt(3), Standardized: true}),
				diag.String("user_local", "test"),
				diag.Field{Key: "std_local", Value: diag.ValInt(2), Standardized: true}),
			filter: filtUser,
			want: map[string]interface{}{
				"user_before": "test",
				"user_local":  "test",
				"user_after":  "test",
			},
			len: 3,
		},

		"standardized filter transitive": {
			in: makeCtx(
				makeCtx(nil, nil, diag.String("user_before", "test"),
					diag.Field{Key: "std_before", Value: diag.ValInt(1), Standardized: true}),
				makeCtx(nil, nil, diag.String("user_after", "test"),
					diag.Field{Key: "std_after", Value: diag.ValInt(3), Standardized: true}),
				diag.String("user_local", "test"),
				diag.Field{Key: "std_local", Value: diag.ValInt(2), Standardized: true}),
			filter: filtStd,
			want: map[string]interface{}{
				"std_before": 1,
				"std_local":  2,
				"std_after":  3,
			},
			len: 3,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := test.filter(test.in)
			assertCtx(t, test.want, ctx)
			requireEqual(t, test.len, ctx.Len())
		})
	}
}

func TestCtxVisitKeyValues(t *testing.T) {
	ctx := makeCtx(nil, nil,
		diag.String("a.b.field1", "test"),
		diag.String("a.b.field2", "test"),
		diag.Int("a.c.field1", 1),
		diag.Int("a.c.field2", 2),
		diag.Int("z.c", 5),
		diag.Int("z.d", 6))

	var v testVisitor
	requireNoError(t, ctx.VisitKeyValues(&v))

	requireEqual(t, map[string]interface{}{
		"a.b.field1": "test",
		"a.b.field2": "test",
		"a.c.field1": 1,
		"a.c.field2": 2,
		"z.c":        5,
		"z.d":        6,
	}, v.Get())
}

func TestCtxVisitStructured(t *testing.T) {
	ctx := makeCtx(nil, nil,
		diag.String("a.b.field1", "test"),
		diag.String("a.b.field2", "test"),
		diag.Int("a.c.field1", 1),
		diag.Int("a.c.field2", 2),
		diag.Int("z.c", 5),
		diag.Int("z.d", 6))

	var v testVisitor
	requireNoError(t, ctx.VisitStructured(&v))

	requireEqual(t, map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"field1": "test",
				"field2": "test",
			},
			"c": map[string]interface{}{
				"field1": 1,
				"field2": 2,
			},
		},
		"z": map[string]interface{}{
			"c": 5,
			"d": 6,
		},
	}, v.Get())
}
