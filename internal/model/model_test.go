package model

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIsSupportedState(t *testing.T) {
	a := assert.New(t)

	a.False(IsSupportedState("foobar"))
	a.False(IsSupportedState("12467"))
	a.True(IsSupportedState(StateLobby))
	a.True(IsSupportedState(strings.ToUpper(StateRunning)))
}

func TestIsSupportedRole(t *testing.T) {
	a := assert.New(t)

	a.False(IsSupportedRole("foobar"))
	a.False(IsSupportedRole("12467"))
	a.True(IsSupportedRole(RoleHost))
	a.True(IsSupportedRole(strings.ToUpper(RolePlayer)))
}
