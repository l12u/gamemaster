package valid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateURL(t *testing.T) {
	a := assert.New(t)

	a.True(ValidateURL("foobar:8000"))
	a.True(ValidateURL("http://foobar:3000/stuff"))
	a.False(ValidateURL("http://foobar:3000#ddd"))
	a.False(ValidateURL("foobar/test"))
	a.False(ValidateURL("://foobar"))
	a.True(ValidateURL("foobar:3000/test"))
	a.True(ValidateURL("http://domain.com/stuff"))
	a.True(ValidateURL("https://postgres:5432"))
	a.True(ValidateURL("http://domain/stuff"))
	a.True(ValidateURL("http://domain/stuff?arg=arg"))
}

func TestValidateSlug(t *testing.T) {
	a := assert.New(t)

	a.False(ValidateSlug("foo+bar"))
	a.False(ValidateSlug("-"))
	a.False(ValidateSlug(""))
	a.False(ValidateSlug("#someStuff"))
	a.True(ValidateSlug("foo-bar"))
	a.False(ValidateSlug("-bar"))
	a.True(ValidateSlug("foo"))
	a.True(ValidateSlug("Some-Stuff-Yea"))
	a.False(ValidateSlug("foo-"))
}
