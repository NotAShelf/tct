package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

var httpClient HttpClient = &http.Client{}

func main() {
	urlPtr := flag.String("url", "http://example.com", "URL to fetch")
	maxRequestsPtr := flag.Int("max", 100, "Maximum number of parallel requests")
	delayPtr := flag.Duration("delay", 0, "Delay between requests")

	flag.Parse()
	url := *urlPtr
	maxRequests := *maxRequestsPtr
	delay := *delayPtr

	// context that listens for the interrupt signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	optimal := 0
	minTime := time.Duration(0)
	for i := 1; i <= maxRequests; i++ {
		// check for cancellation before starting the loop
		if ctx.Err() != nil {
			break
		}
		start := time.Now()
		wg := &sync.WaitGroup{}
		wg.Add(i)
		for j := 0; j < i; j++ {
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
					fmt.Println("Operation cancelled.")
					return
				default:
					time.Sleep(delay)
					makeRequest(httpClient, url)
				}
			}()
		}
		wg.Wait()
		duration := time.Since(start)
		fmt.Printf("Parallel Requests: %d, Time Taken: %s\n", i, duration)
		if minTime == 0 || duration < minTime {
			minTime = duration
			optimal = i
		}
	}
	fmt.Printf("\nOptimal Number of Parallel TCP Requests: %d\n", optimal)
}

func makeRequest(client HttpClient, url string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	_ = body
}
