package environment

import "testing"

func TestName(t *testing.T) {
	p := New_()
	p.WalkDirs()
	p.Orig()
	p.Update()
}
