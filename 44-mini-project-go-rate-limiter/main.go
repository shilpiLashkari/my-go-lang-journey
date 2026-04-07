package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ============================================================
// TOKEN BUCKET LIMITER
// Allows a burst of requests up to 'capacity', then refills
// at 'rate' tokens per second.
// ============================================================

type TokenBucketLimiter struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64 // tokens per second
	lastTime time.Time
}

func NewTokenBucketLimiter(capacity, rate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		tokens:   capacity,
		capacity: capacity,
		rate:     rate,
		lastTime: time.Now(),
	}
}

// Allow returns true if the request is permitted.
func (tb *TokenBucketLimiter) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastTime).Seconds()
	tb.lastTime = now

	// Refill tokens based on elapsed time
	tb.tokens += elapsed * tb.rate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= 1.0 {
		tb.tokens--
		return true
	}
	return false
}

func (tb *TokenBucketLimiter) Status() string {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	return fmt.Sprintf("%.2f / %.2f tokens available", tb.tokens, tb.capacity)
}

// ============================================================
// SLIDING WINDOW LIMITER
// Counts requests in the last 'window' duration.
// Allows up to 'limit' requests per window.
// ============================================================

type SlidingWindowLimiter struct {
	mu         sync.Mutex
	timestamps []time.Time
	limit      int
	window     time.Duration
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:  limit,
		window: window,
	}
}

// Allow returns true if the request is permitted.
func (sw *SlidingWindowLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.window)

	// Evict timestamps outside the window
	valid := sw.timestamps[:0]
	for _, t := range sw.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	sw.timestamps = valid

	if len(sw.timestamps) < sw.limit {
		sw.timestamps = append(sw.timestamps, now)
		return true
	}
	return false
}

func (sw *SlidingWindowLimiter) Status() string {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return fmt.Sprintf("%d / %d requests used in window", len(sw.timestamps), sw.limit)
}

// ============================================================
// HTTP HANDLERS
// ============================================================

func makeTokenBucketHandler(limiter *TokenBucketLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "✅ [Token Bucket] Request allowed! %s\n", limiter.Status())
			log.Printf("[TokenBucket] 200 OK — %s", limiter.Status())
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "🚫 [Token Bucket] Rate limit exceeded! %s\n", limiter.Status())
			log.Printf("[TokenBucket] 429 Too Many Requests — %s", limiter.Status())
		}
	}
}

func makeSlidingWindowHandler(limiter *SlidingWindowLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "✅ [Sliding Window] Request allowed! %s\n", limiter.Status())
			log.Printf("[SlidingWindow] 200 OK — %s", limiter.Status())
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintf(w, "🚫 [Sliding Window] Rate limit exceeded! %s\n", limiter.Status())
			log.Printf("[SlidingWindow] 429 Too Many Requests — %s", limiter.Status())
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "🚀 Go-Rate-Limiter Demo Server")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Endpoints:")
	fmt.Fprintln(w, "  GET /token-bucket    — Token Bucket algorithm (burst-friendly)")
	fmt.Fprintln(w, "  GET /sliding-window  — Sliding Window algorithm (smooth limiting)")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Try hitting each endpoint rapidly to see 429 Too Many Requests!")
}

// ============================================================
// MAIN
// ============================================================

func main() {
	port := flag.String("port", "8080", "HTTP server port")
	tbCapacity := flag.Float64("tb-capacity", 5, "Token Bucket: max burst capacity")
	tbRate := flag.Float64("tb-rate", 1, "Token Bucket: refill rate (tokens/sec)")
	swLimit := flag.Int("sw-limit", 5, "Sliding Window: max requests per window")
	swWindow := flag.Duration("sw-window", 10*time.Second, "Sliding Window: time window duration")
	flag.Parse()

	tbLimiter := NewTokenBucketLimiter(*tbCapacity, *tbRate)
	swLimiter := NewSlidingWindowLimiter(*swLimit, *swWindow)

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/token-bucket", makeTokenBucketHandler(tbLimiter))
	mux.HandleFunc("/sliding-window", makeSlidingWindowHandler(swLimiter))

	fmt.Printf("🚀 Go-Rate-Limiter running on http://localhost:%s\n", *port)
	fmt.Printf("🪣 Token Bucket  : capacity=%.0f, rate=%.1f tokens/sec\n", *tbCapacity, *tbRate)
	fmt.Printf("🪟 Sliding Window: limit=%d requests per %s\n", *swLimit, *swWindow)
	fmt.Println("---------------------------------------------------------")

	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
