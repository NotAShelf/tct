package cmd

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Mock client to count requests
type mockClient struct {
	requestCount int64
}

func (m *mockClient) MakeRequest(url string) {
	atomic.AddInt64(&m.requestCount, 1)
	// Small delay to simulate actual request
	time.Sleep(1 * time.Millisecond)
}

func newMockClient() *mockClient {
	return &mockClient{}
}

func TestGoroutineSpawning(t *testing.T) {
	tests := []struct {
		name           string
		goroutineCount int
		expectedCalls  int64
	}{
		{
			name:           "single goroutine",
			goroutineCount: 1,
			expectedCalls:  1,
		},
		{
			name:           "multiple goroutines",
			goroutineCount: 5,
			expectedCalls:  5,
		},
		{
			name:           "many goroutines",
			goroutineCount: 10,
			expectedCalls:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := newMockClient()
			ctx := context.Background()

			// Simulate the concurrency loop
			wg := &sync.WaitGroup{}
			wg.Add(tt.goroutineCount)

			for range tt.goroutineCount {
				go func() {
					defer wg.Done()
					select {
					case <-ctx.Done():
						return
					default:
						mockClient.MakeRequest("http://test.com")
					}
				}()
			}

			wg.Wait()

			if mockClient.requestCount != tt.expectedCalls {
				t.Errorf("Expected %d requests, got %d", tt.expectedCalls, mockClient.requestCount)
			}
		})
	}
}

func TestConcurrencyWithContext(t *testing.T) {
	mockClient := newMockClient()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel immediately to test cancellation behavior
	cancel()

	wg := &sync.WaitGroup{}
	goroutineCount := 5
	wg.Add(goroutineCount)

	for range goroutineCount {
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				mockClient.MakeRequest("http://test.com")
			}
		}()
	}

	wg.Wait()

	// Since context is cancelled, no requests should be made
	if mockClient.requestCount != 0 {
		t.Errorf("Expected 0 requests due to cancellation, got %d", mockClient.requestCount)
	}
}

func TestConcurrencyTiming(t *testing.T) {
	mockClient := newMockClient()
	ctx := context.Background()
	goroutineCount := 3

	start := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(goroutineCount)

	for range goroutineCount {
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				// Simulate a 5ms request
				time.Sleep(5 * time.Millisecond)
				mockClient.MakeRequest("http://test.com")
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	// With 3 concurrent goroutines, total time should be roughly 5ms, not 15ms
	if duration > 10*time.Millisecond {
		t.Errorf("Concurrent execution took too long: %v, expected around 5ms", duration)
	}

	if mockClient.requestCount != int64(goroutineCount) {
		t.Errorf("Expected %d requests, got %d", goroutineCount, mockClient.requestCount)
	}
}
