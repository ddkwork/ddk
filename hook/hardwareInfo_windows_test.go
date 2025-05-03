package hook

import (
	"testing"

	hardwareinfo "github.com/ddkwork/ddk/hardwareinfo"
	"github.com/ddkwork/golibrary/mylog"
)

func Test_hardware(t *testing.T) {
	// t.Skip()
	h := hardwareinfo.New()
	// if !h.SsdInfo.GetMust() { // todo bug cpu pkg init
	//	return
	// }
	if !h.CpuInfo.Get() {
		return
	}
	if !h.MacInfo.Get() {
		return
	}
	mylog.Struct(h)
}
