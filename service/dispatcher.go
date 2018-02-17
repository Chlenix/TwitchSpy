package service

import (
	"time"
	"fmt"
	"math/rand"
	"sync"
)

const viewerCount = 100

var wg sync.WaitGroup

type Dispatcher struct {
	JobQueue chan int
}

func cclose(channels *[viewerCount]chan int) {
	for _, c := range channels {
		close(c)
	}
}

func Dispatch(out *[viewerCount]chan int) {
	defer wg.Done()
	defer cclose(out)

	videoChunks := make([]int, 100)

	for i := range videoChunks {
		videoChunks[i] = rand.Intn(50000) // we pretend it's a url

		// send the same url to all goroutines
		for j := range out {
			out[j] <- videoChunks[i]
		}

		time.Sleep(40 * time.Millisecond)
	}

	fmt.Println("Dispatcher exited.")
}