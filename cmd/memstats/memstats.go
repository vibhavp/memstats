package main

import (
	"flag"
	"log"
	"time"

	"runtime"

	"github.com/gizak/termui"
	"github.com/vibhavp/memstats/internal/fetch"
)

var (
	target = flag.String("target", "http://localhost:6060", "target link")
	//boxes
	heapAlloc   = termui.NewLineChart()
	heapObjects = termui.NewLineChart()
	gcPause     = termui.NewLineChart()
)

func init() {
	heapAlloc.BorderLabel = "bytes allocated and not yet freed"
	heapObjects.BorderLabel = "total number of allocated objects"
	gcPause.BorderLabel = "GC Pause"

	gcPause.Height = 50
	gcPause.Width = 66

	heapAlloc.Height = 50
	heapAlloc.Width = 66
	heapAlloc.X = 67

	heapObjects.X = 133
	heapObjects.Height = 50
	heapObjects.Width = 66
}

var prevTime float64

func reset(arr *[]float64) {
	if len(*arr) == 80 {
		*arr = []float64{}
	}
}

func render(mstats runtime.MemStats) {
	ns := float64(mstats.PauseNs[(mstats.NumGC+255)%256])
	if prevTime != ns {
		gcPause.Data = append(gcPause.Data, ns)
	}

	reset(&heapAlloc.Data)
	reset(&heapObjects.Data)
	reset(&gcPause.Data)

	heapAlloc.Data = append(heapAlloc.Data, float64(mstats.HeapAlloc))
	heapObjects.Data = append(heapObjects.Data, float64(mstats.HeapObjects))
	prevTime = ns
	//termui.Clear()
	termui.Render(gcPause, heapAlloc, heapObjects)
}

func main() {
	flag.Parse()
	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for {
			mstats, err := fetch.FetchMemStats(*target)
			if err != nil {
				termui.StopLoop()
				log.Fatal(err)
			}
			render(*mstats)
			<-ticker.C
		}
	}()
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		termui.StopLoop()
	})

	termui.Loop()
}
