package ws

import (
	"context"
	"errors"
	"time"
)

// Error messages
var (
	ErrInvalidReconnectionTimeout   = errors.New("invalid reconnection timeout")
	ErrInvalidPongInterval          = errors.New("invalid pong interval")
	ErrInvalidPingInterval          = errors.New("invalid ping interval")
	ErrInvalidPingDeadline          = errors.New("invalid ping deadline")
	ErrInvalidReconnectionBackoff   = errors.New("invalid reconnection backoff")
	ErrInvalidPongHandler           = errors.New("invalid pong handler")
	ErrInvalidCloseHandler          = errors.New("invalid close handler")
	ErrInvalidCompressionLevel      = errors.New("invalid compression level")
	ErrInvalidReadDeadline          = errors.New("invalid read deadline")
	ErrInvalidReadLimit             = errors.New("invalid read limit")
	ErrInvalidWriteDeadline         = errors.New("invalid write deadline")
	ErrInvalidSetCompressionHandler = errors.New("invalid set compression handler")
)

// ClientOpts defines the options for the WebSocket client.
type ClientOpts struct {
	ReconnectionTimeout        time.Duration
	PongInterval               time.Duration
	PingInterval               time.Duration
	PingDeadline               time.Duration
	ReconnectionBackoff        func(ctx context.Context, reconnect func() error, backoff *Backoff) error
	OnConnect                  func(c Client, conn Connection)
	OnDisconnect               func(c Client, conn Connection)
	OnReconnectFailure         func(c Client, conn Connection, err error)
	OnCloseError               func(c Client, conn Connection, err error)
	OnSetCompressionLevelError func(c Client, conn Connection, err error)
	OnSetReadDeadlineError     func(c Client, conn Connection, err error)
	OnSetWriteDeadlineError    func(c Client, conn Connection, err error)
}

// DefaultOpts returns the default client options.
func DefaultOpts() *ClientOpts {
	return &ClientOpts{
		ReconnectionTimeout:        5 * time.Second,
		PongInterval:               60 * time.Second,
		PingInterval:               30 * time.Second,
		PingDeadline:               10 * time.Second,
		ReconnectionBackoff:        DefaultBackoff,
		OnConnect:                  nil,
		OnDisconnect:               nil,
		OnReconnectFailure:         nil,
		OnCloseError:               nil,
		OnSetCompressionLevelError: nil,
		OnSetReadDeadlineError:     nil,
		OnSetWriteDeadlineError:    nil,
	}
}

func checkOpts(opts *ClientOpts) {
	if opts.ReconnectionTimeout <= 0 {
		panic(ErrInvalidReconnectionTimeout)
	}
	if opts.PongInterval <= 0 {
		panic(ErrInvalidPongInterval)
	}
	if opts.PingInterval <= 0 {
		panic(ErrInvalidPingInterval)
	}
	if opts.PingDeadline <= 0 {
		panic(ErrInvalidPingDeadline)
	}
	if opts.ReconnectionBackoff == nil {
		panic(ErrInvalidReconnectionBackoff)
	}
}
