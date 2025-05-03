package drivertool

import (
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr" // todo if build on linux,it need change to cmd

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
)

type (
	helper interface {
		SetService()
		SetManager()
		StartService()
		StopService()
		DeleteService()
		QueryService()
	}
	Interface interface {
		Load(sysPath string)
		Unload()
	}
	Object struct {
		Status       uint32
		service      *mgr.Service
		manager      *mgr.Mgr
		driverPath   string
		DeviceName   string
		dependencies []string
	}
)

func (o *Object) SetDependencies(dependencies []string) {
	o.dependencies = dependencies
}

func NewObject(deviceName, driverPath string) *Object {
	return &Object{
		Status:     0,
		service:    nil,
		manager:    nil,
		driverPath: driverPath,
		DeviceName: deviceName,
	}
}

func New() (d Interface) {
	return NewObject("", "")
}

func (o *Object) Load(sysPath string) {
	o.driverPath = filepath.Join(os.Getenv("SYSTEMROOT"), "system32", "drivers", filepath.Base(sysPath))
	if o.DeviceName == "" {
		o.DeviceName = stream.BaseName(sysPath)
	}
	mylog.Trace("deviceName", o.DeviceName)
	mylog.Trace("driverPath", o.driverPath)
	stream.WriteBinaryFile(o.driverPath, stream.NewBuffer(sysPath).Bytes())
	o.SetManager()
	o.SetService()
	o.StartService()
	mylog.Success("driver load success", o.driverPath)
	o.QueryService()
}

func (o *Object) Unload() {
	o.StopService()
	o.DeleteService()
	mylog.Check(o.manager.Disconnect())
	mylog.Check(o.service.Close())
	mylog.Success("driver unload success", o.driverPath)
	mylog.Check(os.Remove(o.driverPath))
}

func (o *Object) SetService() {
	var e error
	o.service, e = o.manager.OpenService(o.DeviceName)
	if e != nil {
		config := mgr.Config{
			ServiceType:    windows.SERVICE_KERNEL_DRIVER,
			StartType:      mgr.StartManual,
			ErrorControl:   0,
			BinaryPathName: "",
			LoadOrderGroup: "",
			TagId:          0,
			// Dependencies:     o.dependencies,
			ServiceStartName: "",
			DisplayName:      "",
			Password:         "",
			Description:      "",
			SidType:          0,
			DelayedAutoStart: false,
		}
		if o.dependencies != nil {
			config.Dependencies = o.dependencies
		}
		o.service = mylog.Check2(o.manager.CreateService(o.DeviceName, o.driverPath, config))
	}
}

func (o *Object) SetManager() {
	o.manager = mylog.Check2(mgr.Connect())
}

func (o *Object) QueryService() {
	o.Status = mylog.Check2(o.service.Query()).ServiceSpecificExitCode
}

func (o *Object) StopService() {
	status := mylog.Check2(o.service.Control(svc.Stop))
	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			mylog.Check("Timed out waiting for service to stop")
		}
		time.Sleep(300 * time.Millisecond)
		o.QueryService()
		mylog.Trace("Service stopped")
	}
}

func (o *Object) DeleteService() {
	mylog.Check(o.service.Delete())
	mylog.Trace("Service deleted")
	o.QueryService()
}
func (o *Object) StartService() { mylog.Check(o.service.Start()) }
