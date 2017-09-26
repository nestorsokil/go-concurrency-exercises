//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"context"
	"log"
	"time"
)

const nonPremiumTimeSec int64 = 10

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	ctx := context.Background()
	if !u.IsPremium {
		var cancel context.CancelFunc
		allowedSec := nonPremiumTimeSec - u.TimeUsed
		if allowedSec <= 0 {
			log.Printf("User %v has no available time left\n", u.ID)
			return false
		}
		timeout := time.Duration(allowedSec) * time.Second
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	ready := make(chan struct{})
	defer close(ready)

	go runWithSignal(ctx, ready, process, u)
	select {
	case <-ready:
		log.Printf("Process for user %v completed\n", u.ID)
		return true
	case <-ctx.Done():
		log.Printf("Process for user %v timed out\n", u.ID)
		return false
	}
}

func runWithSignal(ctx context.Context, signal chan<- struct{}, process func(), user *User) {
	select {
	case <-ctx.Done():
		log.Println("Context was canceled")
		return
	default:
		start := time.Now()
		process()
		since := time.Since(start)
		user.TimeUsed = user.TimeUsed + int64(since.Seconds())
		log.Printf("Process completed in %v", since)
		signal <- struct{}{}
	}
}

func main() {
	RunMockServer()
}
