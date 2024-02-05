package consumers

import (
	"context"
	"sync"
)

var wg sync.WaitGroup

func StartConsumers(ctx context.Context){
	wg.Add(2)

	go func() {
		defer wg.Done()
		consumeAudioLink(ctx)
	}()

	go func() {
		defer wg.Done()
		consumeTranscript(ctx)
	}()
}