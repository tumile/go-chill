package chill

import (
	"net"
	"net/http"
	"time"
)

type keyGenerator func(r *http.Request) string

type RateLimiter struct {
	window     time.Duration
	max        int
	statusCode int
	message    string
	storage    *storage
	keyGen     keyGenerator
}

func NewRateLimiter(window time.Duration, max int) *RateLimiter {
	lim := RateLimiter{
		window:     window,
		max:        max,
		statusCode: 429,
		message:    "Too many requests, please try again later",
		storage:    newStorage(window),
		keyGen: func(r *http.Request) string {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			uip := net.ParseIP(ip)
			if uip == nil {
				return ""
			}
			return uip.String()
		},
	}
	return &lim
}

func (lim *RateLimiter) RateLimit(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key := lim.keyGen(r)
		cur := lim.storage.increment(key)

		if cur > lim.max {
			w.Header().Set("Retry-After", (lim.window / time.Second).String())
			http.Error(w, lim.message, lim.statusCode)
			return
		}

		next(w, r)
	}
}

func (lim *RateLimiter) SetStatusCode(statusCode int) {
	lim.statusCode = statusCode
}

func (lim *RateLimiter) SetMessage(message string) {
	lim.message = message
}

func (lim *RateLimiter) SetKeyGenerator(keyGen keyGenerator) {
	lim.keyGen = keyGen
}
