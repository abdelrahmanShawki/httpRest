package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type user struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var mutex sync.RWMutex
	var users = make(map[string]*user)

	go func() {
		for {
			time.Sleep(60 * time.Minute)
			mutex.RLock()
			for ip, user := range users {
				if time.Since(user.lastSeen) > 30*time.Minute {
					delete(users, ip)
				}
			}
			mutex.RUnlock()
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.config.limiter.enabled {
			next.ServeHTTP(w, r)
			return
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		mutex.Lock()
		if _, found := users[ip]; !found {
			users[ip] = &user{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
		}

		users[ip].lastSeen = time.Now()
		currentLimit := users[ip].limiter

		if !currentLimit.Allow() {
			mutex.Unlock()
			app.rateLimitResponse(w, r)
			return
		}

		mutex.Unlock()
		next.ServeHTTP(w, r)
	})

}
