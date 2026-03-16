package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Item represents a single cached value with an expiration time.
type Item struct {
	Value      interface{}
	Expiration int64 // Unix timestamp in nanoseconds
}

// Cache holds our data map and a Read-Write Mutex to prevent data races.
type Cache struct {
	items map[string]Item
	mu    sync.RWMutex // Note: RWMutex allows multiple simultaneous readers!
}

// NewCache initializes a new cache and starts the background janitor.
func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]Item),
	}
	
	// Start the background janitor
	go c.startJanitor(cleanupInterval)
	return c
}

// Set adds an item to the cache, replacing it if it already exists.
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	// We MUST use a full Lock for writing. No one else can read or write right now.
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
}

// Get retrieves an item from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	// We use RLock! Multiple goroutines can call Get() simultaneously safely.
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock() // Unlock quickly before doing expiration logic

	if !found {
		return nil, false
	}

	// Check if the item has expired (passive deletion)
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		// It expired! Delete it and pretend it doesn't exist.
		c.Delete(key)
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache.
func (c *Cache) Delete(key string) {
	// Full lock for writing
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// startJanitor runs forever, periodically checking for and removing expired items (active deletion).
func (c *Cache) startJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		c.deleteExpired()
	}
}

func (c *Cache) deleteExpired() {
	now := time.Now().UnixNano()
	
	// We must lock the map while checking and deleting items
	c.mu.Lock()
	defer c.mu.Unlock()
	
	deletedCount := 0
	for key, item := range c.items {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.items, key)
			deletedCount++
		}
	}
	
	if deletedCount > 0 {
		fmt.Printf("🧹 [JANITOR] Swept away %d expired items from memory.\n", deletedCount)
	}
}


func main() {
	fmt.Println("🧠 Starting Go-Cache Simulation")
	fmt.Println("==================================================")

	// Create a cache with a janitor that runs every 2 seconds
	cache := NewCache(2 * time.Second)

	// --- 1. Basic Set and Get ---
	fmt.Println("\n1. Basic Functionality:")
	cache.Set("hero", "Batman", 0) // No expiration
	val, ok := cache.Get("hero")
	fmt.Printf("Get 'hero': %v (Found: %v)\n", val, ok)

	// --- 2. TTL Expiration ---
	fmt.Println("\n2. Time-To-Live (TTL) Simulation:")
	fmt.Println("Setting 'secret_code' to '12345' for 3 seconds...")
	cache.Set("secret_code", "12345", 3 * time.Second)
	
	val, ok = cache.Get("secret_code")
	fmt.Printf("-> Immediate Get: %v (Found: %v)\n", val, ok)
	
	fmt.Println("-> Sleeping for 4 seconds...")
	time.Sleep(4 * time.Second)
	
	val, ok = cache.Get("secret_code")
	fmt.Printf("-> Get after sleep: %v (Found: %v)\n", val, ok)

	// --- 3. Concurrency Stress Test ---
	fmt.Println("\n3. Thread-Safety Stress Test:")
	fmt.Println("Spawning 100 concurrent workers to bombard the cache...")
	
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			key := fmt.Sprintf("key_%d", rand.Intn(10))
			
			// 50% chance to Read, 50% chance to Write
			if rand.Float32() < 0.5 {
				cache.Set(key, workerID, 5 * time.Second)
			} else {
				cache.Get(key)
			}
		}(i)
	}
	
	wg.Wait()
	fmt.Println("✅ 100 Goroutines completed without triggering a data race panic!")
	fmt.Println("==================================================")
	
	// Let the janitor run one last time to clean up the stress test mess
	fmt.Println("Waiting for Janitor cleanup...")
	time.Sleep(6 * time.Second)
}
