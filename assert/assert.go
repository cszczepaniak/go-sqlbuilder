package assert

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal compares expected and actual with go-cmp and fails if they differ.
func Equal[T any](tb testing.TB, expected, actual T) {
	tb.Helper()
	if cmp.Equal(expected, actual) {
		return
	}
	tb.Fatalf("not equal:\ndiff: %s", cmp.Diff(expected, actual))
}

// Nil fails if value is not nil.
func Nil(tb testing.TB, value any) {
	tb.Helper()
	if value == nil {
		return
	}
	if rv := reflect.ValueOf(value); rv.Kind() == reflect.Ptr && rv.IsNil() {
		return
	}
	tb.Fatalf("expected nil, got %v", value)
}

// NotNil fails if value is nil.
func NotNil(tb testing.TB, value any) {
	tb.Helper()
	if value == nil {
		tb.Fatal("expected non-nil value")
	}
	if rv := reflect.ValueOf(value); rv.Kind() == reflect.Ptr && rv.IsNil() {
		tb.Fatal("expected non-nil value")
	}
}

// NoError fails if err is not nil.
func NoError(tb testing.TB, err error) {
	tb.Helper()
	if err == nil {
		return
	}
	tb.Fatalf("unexpected error: %v", err)
}

// Error fails if err is nil.
func Error(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		return
	}
	tb.Fatal("expected error, got nil")
}

// LessOrEqual fails if a > b.
func LessOrEqual(tb testing.TB, a, b int) {
	tb.Helper()
	if a <= b {
		return
	}
	tb.Fatalf("expected %d <= %d", a, b)
}
