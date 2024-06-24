package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testType struct{}

func (t *testType) aFunc() {}

func (t testType) anotherFunc() {}

func TestGetFunctionName(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		f    any
		name string
	}{
		// Basic function
		{
			f:    GetFunctionName,
			name: "internal.GetFunctionName",
		},

		// Basic function in another package
		{
			f:    testing.Benchmark,
			name: "testing.Benchmark",
		},

		// Receiver functions
		{
			f:    (*testType)(nil).aFunc,
			name: "internal.(*testType).aFunc",
		},
		{
			f:    testType{}.anotherFunc,
			name: "internal.testType.anotherFunc",
		},
	}

	for i, testCase := range testCases {
		assert.Equal(t, testCase.name, GetFunctionName(testCase.f), "mismatched function name for id %d", i)
	}
}

func TestCurrentFunc(t *testing.T) {
	t.Parallel()

	var innerCurrentFunc = func(skip int) string {
		return CurrentFunc(skip)
	}

	assert.Equal(t, "internal.TestCurrentFunc", CurrentFunc(1))
	assert.Equal(t, "internal.TestCurrentFunc.func1", innerCurrentFunc(1))
	assert.Equal(t, "internal.TestCurrentFunc", innerCurrentFunc(2))

	// This actually happens at 5, but we pick an arbitrarily large number in
	// case testing internals change.
	assert.Equal(t, "unknown", innerCurrentFunc(42))
}
