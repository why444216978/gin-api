package logging

import "testing"

func TestAddTag(t *testing.T) {
	l := &LogHeader{}
	l.AddTag("a", "b", "b", "c")
	t.Log(l)
	l.SetTag("c", "e", "e")
	t.Log(l)
}
