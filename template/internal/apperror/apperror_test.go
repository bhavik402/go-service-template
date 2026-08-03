package apperror

import (
	"errors"
	"testing"
)

func TestError_UnwrapAndIs(t *testing.T) {
	cause := errors.New("connection refused")
	err := Internal("database error", cause)

	if !errors.Is(err, cause) {
		t.Error("expected errors.Is to find the wrapped cause")
	}

	var appErr *Error
	if !errors.As(err, &appErr) {
		t.Fatal("expected errors.As to match *Error")
	}
	if appErr.Code != CodeInternal {
		t.Errorf("expected code %s, got %s", CodeInternal, appErr.Code)
	}
}

func TestNotFound_HasNoUnderlyingCause(t *testing.T) {
	err := NotFound("user not found")

	if err.Unwrap() != nil {
		t.Error("expected NotFound to have no wrapped cause")
	}
	if err.Error() != "user not found" {
		t.Errorf("unexpected message: %s", err.Error())
	}
}
