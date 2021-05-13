package log

import "testing"

func TestLevel(t *testing.T) {
	for i := int(MinLevel); i <= int(MaxLevel); i++ {
		lvl := Level(i)
		SetLevel(lvl)
		if GetLevel() != lvl {
			t.Errorf("bad level, want: %v, got: %v", lvl, GetLevel())
		}
	}
}
