package optstr

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseString(t *testing.T) {
	l, e := ParseString(`aaa1=111 bbb="世界"`)
	assert.NoError(t, e)
	t.Log(l, e)
	l, e = ParseString(`aaa1=111 bbb="世界\"`)
	assert.Error(t, e)
	t.Log(l, e)
	l, e = ParseString(`aaa1=111 bbb= `)
	assert.NoError(t, e)
	t.Log(l, e)
	l, e = ParseString(`aaa1=111 bbb=`)
	assert.NoError(t, e)
	t.Log(l, e)
}

func Test_ParseString2(t *testing.T) {
	l, e := ParseString(`a= `)
	assert.NoError(t, e)
	t.Log(l, e)
}
