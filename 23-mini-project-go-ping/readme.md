# Day 23: Go-Ping - Concurrent Health Checker 🏓🌐

## Overview

For Day 23, I decided to build **Go-Ping**, a straightforward but powerful concurrent website uptime monitor. Monitoring the health of external services is a crucial part of backend development. This project takes a list of URLs and checks them all simultaneously, drastically reducing the total time it takes to sweep an entire list of sites compared to a synchronous loop.

## Features

- **Concurrent Execution**: Spawns an individual goroutine for every URL, meaning 100 checks take roughly the exact same time as 1 check!
- **Channel Aggregation**: Safely funnels the results from dozens of scattered goroutines back into a single thread for organized printing.
- **Strict Timeouts**: Uses a customized `http.Client` with a timeout to guarantee the program never hangs infinitely on a broken server block. 
- **Detailed Output**: Reports HTTP status codes, whether the site is UP or DOWN, and the exact millisecond time it took to respond.

## How to Run

1. Navigate to the directory:
   ```bash
   cd 23-mini-project-go-ping
   ```
2. Run the program:
   ```bash
   go run main.go
   ```
3. Watch as it immediately spits out the varying health statuses of multiple simulated websites in parallel!

## Learning Reflection

- **Scatter/Gather Pattern**: Practiced the classic Go pattern of "scattering" identical work across multiple goroutines, and "gathering" the replies back via a shared `chan`. 
- **Error Handling vs. Panics**: Reinforced the concept of treating networking errors (like DNS lookup failures or timeouts) as valid data to be returned via the struct, rather than letting the application crash.
- **HTTP Client**: Learned the importance of explicitly setting `Timeout: 5 * time.Second` on `http.Client`. The default Go HTTP client has NO timeout, which is a major production hazard!

---
*Ping... Pong! Master of concurrent networking. 🐹*
