package drivertool

import (
	"testing"
)

func TestLoadSys(t *testing.T) {
	path := "sysDemo.sys"
	d := New()
	d.Load(path)
	d.Unload()
}
