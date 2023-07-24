package backupper_go

import (
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

type Backupper struct {
	state int32
	Run   func()
}

const (
	NoRunning     = 0
	Running       = 1
	RunningBackup = 2
)

func (b *Backupper) DoRun() {
	for {
		if b.state == NoRunning {
			if atomic.CompareAndSwapInt32(&b.state, NoRunning, Running) {
				b.Run()
				for {
					if b.state == RunningBackup {
						atomic.StoreInt32(&b.state, Running)
						b.Run()
					} else {
						break
					}
				}
				if atomic.CompareAndSwapInt32(&b.state, Running, NoRunning) {
					break
				}
				atomic.StoreInt32(&b.state, NoRunning)
			}
		} else if b.state == Running {
			if atomic.CompareAndSwapInt32(&b.state, Running, RunningBackup) {
				break
			}
		} else if b.state == RunningBackup {
			break
		}
	}
}

func getGoId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	field := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, _ := strconv.Atoi(field)
	return id
}
