package chill

import "time"

type storage struct {
	hits map[string]int
}

func newStorage(window time.Duration) *storage {
	s := storage{
		hits: make(map[string]int),
	}
	ticker := time.NewTicker(window)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.reset()
			}
		}
	}()
	return &s
}

func (s *storage) reset() {
	s.hits = map[string]int{}
}

func (s *storage) increment(key string) int {
	s.hits[key]++
	return s.hits[key]
}

func (s *storage) decrement(key string) {
	s.hits[key]--
}
