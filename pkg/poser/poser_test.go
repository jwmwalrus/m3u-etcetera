package poser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type posign struct {
	msg      string
	position int
	ignore   bool
}

func (p *posign) GetPosition() int {
	return p.position
}

func (p *posign) SetPosition(pos int) {
	p.position = pos
}

func (p *posign) GetIgnore() bool {
	return p.ignore
}

func (p *posign) SetIgnore(ignore bool) {
	p.ignore = ignore
}

func TestAppendTo(t *testing.T) {
	s := []*posign{
		{"", 1, false},
		{"", 2, false},
		{"", 3, false},
		{"", 4, false},
	}

	s = AppendTo(s, &posign{})

	assert.Equal(t, len(s), 5)
	assert.Equal(t, s[4].GetPosition(), 5)
}

func TestDeleteAt(t *testing.T) {
	s := []*posign{
		{"", 1, false},
		{"", 2, false},
		{"", 3, false},
		{"", 4, false},
	}

	s, e := DeleteAt(s, 3)

	assert.Equal(t, len(s), 3)
	assert.Equal(t, s[2].GetPosition(), 3)
	assert.Equal(t, e.GetPosition(), 3)
}

func TestInsertInto(t *testing.T) {
	s := []*posign{
		{"one", 1, false},
		{"two", 2, false},
		{"four", 3, false},
		{"five", 4, false},
	}

	s = InsertInto(s, 3, &posign{msg: "three"})

	assert.Equal(t, len(s), 5)
	assert.Equal(t, s[3].GetPosition(), 4)
	assert.Equal(t, s[2].msg, "three")
}

func TestMoveTo(t *testing.T) {
	s := []*posign{
		{"one", 1, false},
		{"four", 2, false},
		{"three", 3, false},
		{"two", 4, false},
	}

	s = MoveTo(s, 2, 4)

	assert.Equal(t, len(s), 4)
	assert.Equal(t, s[1].GetPosition(), 2)
	assert.Equal(t, s[1].msg, "two")

	s = MoveTo(s, 4, 3)

	assert.Equal(t, len(s), 4)
	assert.Equal(t, s[2].msg, "three")
	assert.Equal(t, s[3].msg, "four")
}

func TestPop(t *testing.T) {
	s := []*posign{
		{"one", 1, false},
		{"two", 2, false},
		{"three", 3, false},
		{"four", 4, false},
	}

	s, e := Pop(s)

	assert.Equal(t, len(s), 3)
	assert.Equal(t, s[0].GetPosition(), 1)
	assert.Equal(t, e.msg, "one")
}

func TestPrependItem(t *testing.T) {
	s := []*posign{
		{"one", 1, false},
		{"two", 2, false},
		{"three", 3, false},
		{"four", 4, false},
	}

	s = PrependItem(s, &posign{msg: "zero"})

	assert.Equal(t, len(s), 5)
	assert.Equal(t, s[4].GetPosition(), 5)
	assert.Equal(t, s[0].msg, "zero")
}
