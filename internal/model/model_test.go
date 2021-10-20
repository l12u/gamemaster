package model

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIsSupportedState(t *testing.T) {
	a := assert.New(t)

	a.False(IsSupportedRole("foobar"))
	a.False(IsSupportedRole("12467"))
	a.True(IsSupportedRole(StateLobby))
	a.True(IsSupportedRole(strings.ToUpper(StateRunning)))
}

func TestIsSupportedRole(t *testing.T) {
	a := assert.New(t)

	a.False(IsSupportedRole("foobar"))
	a.False(IsSupportedRole("12467"))
	a.True(IsSupportedRole(RoleHost))
	a.True(IsSupportedRole(strings.ToUpper(RolePlayer)))
}
