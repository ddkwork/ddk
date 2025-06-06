package taskmanager

import (
	"github.com/ddkwork/golibrary/std/mylog"
	"github.com/ddkwork/golibrary/std/stream/datasize"
	"github.com/ddkwork/ux"
	"github.com/shirou/gopsutil/v3/process"
)

type Task struct {
	Name    string
	CPU     float64
	RAM     datasize.Size
	RAMPct  float32 `label:"RAM %" format:"%.3g%%"`
	Threads int32
	User    string
	PID     int32
	// todo add hex pid, net walk, inject etc
}

func getTasks(parent *ux.Node[Task]) {
	ps := mylog.Check2(process.Processes())
	for _, p := range ps {
		// todo query Bs path abd read icon
		t := Task{
			Name:    mylog.Check2Ignore(p.Name()),
			CPU:     mylog.Check2Ignore(p.CPUPercent()),
			RAM:     0,
			RAMPct:  mylog.Check2Ignore(p.MemoryPercent()),
			Threads: mylog.Check2Ignore(p.NumThreads()),
			User:    mylog.Check2Ignore(p.Username()),
			PID:     p.Pid,
		}
		mi := mylog.Check2Ignore(p.MemoryInfo())
		if mi != nil {
			t.RAM = datasize.Size(mi.RSS)
		}

		children, e := p.Children()
		mylog.CheckIgnore(e)
		if e != nil || len(children) == 0 || children == nil {
			parent.AddChildByData(t)
			continue
		}
		if len(children) > 0 {
			container := ux.NewContainerNode(mylog.Check2Ignore(p.Name()), t)
			parent.AddChild(container)
			for _, child := range children {
				data := Task{
					Name:    mylog.Check2Ignore(child.Name()),
					CPU:     mylog.Check2Ignore(child.CPUPercent()),
					RAM:     0,
					RAMPct:  mylog.Check2Ignore(child.MemoryPercent()),
					Threads: mylog.Check2Ignore(child.NumThreads()),
					User:    mylog.Check2Ignore(child.Username()),
					PID:     child.Pid,
				}
				mi := mylog.Check2Ignore(child.MemoryInfo())
				if mi != nil {
					data.RAM = datasize.Size(mi.RSS)
				}
				container.AddChildByData(data)
			}
		}
	}
}
