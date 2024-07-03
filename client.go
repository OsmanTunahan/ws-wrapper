package ws

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Client interface {
	Dial(urlStr string, requestHeader http.Header) (*ClientImpl, *http.Response, error)
	DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*ClientImpl, *http.Response, error)
	Reconnect(ctx context.Context) <-chan bool
	Connected() bool
	Closed() bool
	Close() error
	Connection() *websocket.Conn
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type ClientImpl struct {
	urlStr        string
	requestHeader http.Header
	dialer        *websocket.Dialer
	opts          *ClientOpts

	enableWriteCompression *bool
	compressionLevel       *int
	readDeadline           *time.Time
	readLimit              *int64
	writeDeadline          *time.Time
	handlePing             func(string) error
	handleClose            func(int, string) error

	conn       atomic.Value
	handlePong atomic.Value
	pong       atomic.Int64
	connected  atomic.Bool
	closed     atomic.Bool
	close      sync.Once
	mu         sync.Mutex
	reconnect  chan chan<- bool
	done       chan struct{}
}

func NewClient(urlStr string, requestHeader http.Header, opts *ClientOpts) *ClientImpl {
	return &ClientImpl{
		urlStr:        urlStr,
		requestHeader: requestHeader,
		dialer:        websocket.DefaultDialer,
		opts:          opts,
		done:          make(chan struct{}),
	}
}

func (c *ClientImpl) store(conn *websocket.Conn) {
	if c.enableWriteCompression != nil {
		conn.EnableWriteCompression(*c.enableWriteCompression)
	}
	if c.compressionLevel != nil {
		if err := conn.SetCompressionLevel(*c.compressionLevel); err != nil && c.opts.OnSetCompressionLevelError != nil {
			c.opts.OnSetCompressionLevelError(c, conn, err)
		}
	}
	if c.readDeadline != nil {
		if err := conn.SetReadDeadline(*c.readDeadline); err != nil && c.opts.OnSetReadDeadlineError != nil {
			c.opts.OnSetReadDeadlineError(c, conn, err)
		}
	}
	if c.readLimit != nil {
		conn.SetReadLimit(*c.readLimit)
	}
	if c.writeDeadline != nil {
		if err := conn.SetWriteDeadline(*c.writeDeadline); err != nil && c.opts.OnSetWriteDeadlineError != nil {
			c.opts.OnSetWriteDeadlineError(c, conn, err)
		}
	}

	conn.SetCloseHandler(c.handleClose)
	conn.SetPingHandler(c.handlePing)

	conn.SetPongHandler(func(appData string) error {
		c.pong.Store(time.Now().UnixNano())
		if err := c.handlePong.Load().(func(string) error)(appData); err != nil {
			return err
		}
		return nil
	})

	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.Store(&connectionImpl{conn: conn, done: make(chan struct{})})
}

func (c *ClientImpl) Connection() *websocket.Conn {
	return c.conn.Load().(*connectionImpl).conn
}

func (c *ClientImpl) Reconnect(ctx context.Context) <-chan bool {
	reconnect := make(chan bool, 1)
	select {
	case c.reconnect <- reconnect:
		return reconnect
	case <-c.done:
	case <-ctx.Done():
	}
	reconnect <- false
	return reconnect
}

func (c *ClientImpl) connect() bool {
	c.connected.Store(false)
	connection := c.conn.Load().(*connectionImpl)
	if err := connection.conn.Close(); err != nil && c.opts.OnCloseError != nil {
		c.opts.OnCloseError(c, c.Connection(), err)
	}
	if c.opts.OnDisconnect != nil {
		c.opts.OnDisconnect(c, connection.conn)
	}
	defer close(connection.done)
	err := c.opts.ReconnectionBackoff(context.Background(), func() error {
		ctx, cancel := context.WithTimeout(context.Background(), c.opts.ReconnectionTimeout)
		defer cancel()
		conn, _, err := c.dialer.DialContext(ctx, c.urlStr, c.requestHeader)
		if err != nil {
			return err
		}
		c.store(conn)
		c.connected.Store(true)
		c.pong.Store(time.Now().UnixNano())
		if c.opts.OnConnect != nil {
			c.opts.OnConnect(c, conn)
		}
		return nil
	}, nil)
	if err != nil {
		if err := c.Close(); err != nil && c.opts.OnCloseError != nil {
			c.opts.OnCloseError(c, c.Connection(), err)
		}
		if c.opts.OnReconnectFailure != nil {
			c.opts.OnReconnectFailure(c, c.Connection(), err)
		}
		return false
	}
	return true
}

func (c *ClientImpl) run() {
	ticker := time.NewTicker(c.opts.PingInterval)
	defer ticker.Stop()
	for {
		select {
		case reconnect := <-c.reconnect:
			reconnect <- c.connect()
		case <-ticker.C:
			if time.Since(time.Unix(0, c.pong.Load())) > c.opts.PongInterval {
				c.connect()
			}
			if err := c.Connection().WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(c.opts.PingDeadline)); err != nil {
				c.connect()
			}
		case <-c.done:
			return
		}
	}
}

func (c *ClientImpl) Dial(urlStr string, requestHeader http.Header) (*ClientImpl, *http.Response, error) {
	return c.DialContext(context.Background(), urlStr, requestHeader)
}

func (c *ClientImpl) DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (*ClientImpl, *http.Response, error) {
	if c.Connected() {
		return c, nil, nil
	}

	c.opts = DefaultOpts()

	conn, resp, err := c.dialer.DialContext(ctx, urlStr, requestHeader)
	if err != nil {
		return nil, resp, err
	}
	c.store(conn)
	c.connected.Store(true)
	c.pong.Store(time.Now().UnixNano())
	if c.opts.OnConnect != nil {
		c.opts.OnConnect(c, conn)
	}
	go c.run()
	return c, resp, nil
}

func (c *ClientImpl) Connected() bool {
	return c.connected.Load()
}

func (c *ClientImpl) Closed() bool {
	return c.closed.Load()
}

func (c *ClientImpl) Close() error {
	if !c.Closed() {
		return nil
	}
	c.closed.Store(true)
	c.close.Do(func() { close(c.done) })
	return c.Connection().Close()
}

func (c *ClientImpl) ReadMessage() (messageType int, p []byte, err error) {
	return c.Connection().ReadMessage()
}

func (c *ClientImpl) WriteMessage(messageType int, data []byte) error {
	return c.Connection().WriteMessage(messageType, data)
}
