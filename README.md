## Overview

Resilient WebSocket wrapper for the Gorilla WebSocket client with auto-reconnect on disconnect and thread safety

## Features

- **Flexible Options**: Configure WebSocket connection options such as reconnection timeout, ping/pong intervals, and handlers for various events.
- **Automatic Reconnection**: Handles WebSocket connection drops and attempts automatic reconnection based on configured backoff strategies.
- **Custom Event Handlers**: Define custom handlers for connection events like connect, disconnect, close errors, etc.
- **Error Handling**: Comprehensive error handling for WebSocket connection errors and configuration issues.
- **Concurrency**: Utilizes Go's concurrency features for handling WebSocket connections efficiently.

## Usage

### Installation

Install the library using Go modules:

```bash
go get github.com/OsmanTunahan/ws-wrapper
```

### Example

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OsmanTunahan/ws-wrapper/ws"
)

func main() {
	url := "ws://example.com/ws"
	headers := http.Header{"Authorization": []string{"Bearer your_token"}}

	opts := ws.DefaultOpts()
	opts.ReconnectionTimeout = 10 * time.Second

	client := ws.NewClient(url, headers, opts)
	defer client.Close()

	// Handle WebSocket events
	client.OnConnect = func(c ws.Client, conn ws.Connection) {
		fmt.Println("Connected to WebSocket server")
	}

	// Start the WebSocket client
	if err := client.Connect(); err != nil {
		fmt.Printf("Error connecting to WebSocket: %v\n", err)
		return
	}

	// Perform WebSocket operations
	// e.g., client.Send([]byte("Hello, server"))

	// Keep the main goroutine alive
	select {}
}
```