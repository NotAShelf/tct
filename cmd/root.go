package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"notashelf.dev/tct/internal/tct"
)

var (
	url         string
	maxRequests int
	delay       time.Duration
	Version     string // will be set by main package
)

var rootCmd = &cobra.Command{
	Use:   "tct",
	Short: "TCP Connection Timer - find optimal parallel request count",
	Long:  `A tool to measure and find the optimal number of parallel TCP requests for a given URL.`,
	Run:   run,
}

func Execute() {
	if Version != "" {
		rootCmd.Version = Version
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "http://example.com", "URL to fetch")
	rootCmd.Flags().IntVarP(&maxRequests, "max", "m", 100, "Maximum number of parallel requests")
	rootCmd.Flags().DurationVarP(&delay, "delay", "d", 0, "Delay between requests")
}

func run(cmd *cobra.Command, args []string) {
	client := tct.NewClient()

	// Context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	optimal := 0
	minTime := time.Duration(0)

	for i := 1; i <= maxRequests; i++ {
		// Check for cancellation before starting the loop
		if ctx.Err() != nil {
			break
		}

		start := time.Now()
		wg := &sync.WaitGroup{}
		wg.Add(i)

		for range i {
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
					fmt.Println("Operation cancelled.")
					return
				default:
					time.Sleep(delay)
					client.MakeRequest(url)
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
