package apperr

import (
	"errors"
	"testing"
)

func TestAppError_Unwrap(t *testing.T) {
	wrapped := Wrap(CodeNotFound, ErrNotFound)

	if !errors.Is(wrapped, ErrNotFound) {
		t.Fatalf("Wrap 之后 errors.Is 应该能找到根因 —— 说明 Unwrap() 没生效")
	}
}
