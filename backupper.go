package backupper_go

import (
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

type Backupper struct {
	state int32
	Run func()
}

const (
	NoRunning     = 0
	Running       = 1
	RunningBackup = 2
)

func (b *Backupper) DoRun()  {
	log.Printf("%d enter 0", getGoId())
	for {
		if b.state == NoRunning {
			log.Printf("%d enter 1", getGoId())
			if atomic.CompareAndSwapInt32(&b.state, NoRunning, Running) {
				log.Printf("%d enter 2", getGoId())
				b.Run()
				for {
					if b.state == RunningBackup {
						log.Printf("%d enter 3", getGoId())
						atomic.StoreInt32(&b.state, Running)
						b.Run()
					} else {
						log.Printf("%d enter 4", getGoId())
						break
					}
				}
				log.Printf("%d enter 5", getGoId())
				if atomic.CompareAndSwapInt32(&b.state, Running, NoRunning) {
					log.Printf("%d enter 6", getGoId())
					break
				}
				log.Printf("%d enter 7", getGoId())
				atomic.StoreInt32(&b.state, NoRunning)
			}
		} else if b.state == Running {
			log.Printf("%d enter 8", getGoId())
			if atomic.CompareAndSwapInt32(&b.state, Running, RunningBackup) {
				log.Printf("%d enter 9", getGoId())
				break
			}
			log.Printf("%d enter 10", getGoId())
		} else if b.state == RunningBackup {
			log.Printf("%d enter 11", getGoId())
			break
		} else {
			log.Printf("%d enter 12", getGoId())
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