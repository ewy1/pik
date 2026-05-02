package main

import (
	"github.com/ewy1/pik/spool"
	"sync"
)

// ComponentList is a list wrapper type which handles operation modes of the program
type ComponentList[T any] []T

// RunAsync checks components one by one and triggers the fire method in a goroutine.
// the method will return when the waitgroup is done (all initializers are finished)
func (c ComponentList[T]) RunAsync(fire func(T) error) {
	wg := sync.WaitGroup{}
	for _, i := range c {
		wg.Go(func() {
			err := fire(i)
			if err != nil {
				_, _ = spool.Warn("%v\n", err)
			}
		})
	}
	wg.Wait()
}

// RunSync checks components one by one and fires them synchronously.
// important when the order of init matters (for example, paths needs to go before cache)
func (c ComponentList[T]) RunSync(fire func(T) error) {
	for _, i := range c {
		err := fire(i)
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
		}
	}
}
