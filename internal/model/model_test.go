package model

import (
	"strings"
	"testing"
)

func TestIsSupportedState(t *testing.T) {
	s0 := "foobar"
	s1 := "12467"
	s2 := StateLobby
	s3 := strings.ToUpper(StateRunning)

	if IsSupportedState(s0) {
		t.Errorf("%s should not be a supported state", s0)
	}
	if IsSupportedState(s1) {
		t.Errorf("%s should not be a supported state", s1)
	}
	if !IsSupportedState(s2) {
		t.Errorf("%s should be a supported state", s2)
	}
	if !IsSupportedState(s3) {
		t.Errorf("%s should be a supported state", s3)
	}
}
