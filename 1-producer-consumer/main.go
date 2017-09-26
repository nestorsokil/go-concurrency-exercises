//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream, done chan<- struct{}) chan *Tweet {
	res := make(chan *Tweet)
	go func() {
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				done <- struct{}{}
				return
			}
			res <- tweet
		}
	}()
	return res
}

func consumer(prod chan *Tweet, done <-chan struct{}) {
	for {
		select {
		case t := <-prod:
			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		case <-done:
			return
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	done := make(chan struct{})
	defer close(done)
	tweets := producer(stream, done)
	defer close(tweets)

	nroutines := 5
	var signals []chan struct{}
	for i := 0; i < nroutines; i++ {
		closeSignal := make(chan struct{})
		signals = append(signals, closeSignal)
		go consumer(tweets, closeSignal)
	}

	<-done
	for _, s := range signals {
		close(s)
	}

	fmt.Printf("Process took %s\n", time.Since(start))
}
