package main

import (
	"github.com/ewy1/pik/spool"
	"sync"
)

type ComponentList[T any] []T

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

func (c ComponentList[T]) RunSync(fire func(T) error) {
	for _, i := range c {
		err := fire(i)
		if err != nil {
			_, _ = spool.Warn("%v\n", err)
		}
	}
}
