package backupper_go

import (
	"log"
	"sync"
	"testing"
	"time"
)

var flag = 0

func add() {
	time.Sleep(0)
	log.Printf("%d === %d", getGoId(), flag)
	flag++
}

func TestBackupper_DoRun(t *testing.T) {
	b := Backupper{
		state: 0,
		Run:   add,
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(backupper *Backupper) {
			backupper.DoRun()
			wg.Done()
		}(&b)
	}
	wg.Wait()
}
