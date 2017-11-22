package reload

import (
	"log"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	go func() {
		err := Do(log.Printf)
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)

	// TODO: maybe write some meaningful tests?
}
