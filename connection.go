package ws

import (
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type Connection interface {
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	NextReader() (messageType int, r io.Reader, err error)
	NextWriter(messageType int) (io.WriteCloser, error)
	WriteControl(messageType int, data []byte, deadline time.Time) error
	EnableWriteCompression(enable bool)
	SetCompressionLevel(level int) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetReadLimit(limit int64)
	Subprotocol() string
	UnderlyingConn() net.Conn
	SetPingHandler(h func(appData string) error)
	SetPongHandler(h func(appData string) error)
	SetCloseHandler(h func(code int, text string) error)
	CloseHandler() func(code int, text string) error
	PingHandler() func(appData string) error
	PongHandler() func(appData string) error
}

type connectionImpl struct {
	conn *websocket.Conn
	done chan struct{}
}

func (c *connectionImpl) ReadJSON(v interface{}) error {
	return c.conn.ReadJSON(v)
}

func (c *connectionImpl) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

func (c *connectionImpl) Close() error {
	return c.conn.Close()
}

func (c *connectionImpl) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *connectionImpl) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *connectionImpl) NextReader() (messageType int, r io.Reader, err error) {
	return c.conn.NextReader()
}

func (c *connectionImpl) NextWriter(messageType int) (io.WriteCloser, error) {
	return c.conn.NextWriter(messageType)
}

func (c *connectionImpl) WriteControl(messageType int, data []byte, deadline time.Time) error {
	return c.conn.WriteControl(messageType, data, deadline)
}

func (c *connectionImpl) EnableWriteCompression(enable bool) {
	c.conn.EnableWriteCompression(enable)
}

func (c *connectionImpl) SetCompressionLevel(level int) error {
	return c.conn.SetCompressionLevel(level)
}

func (c *connectionImpl) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *connectionImpl) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *connectionImpl) SetReadLimit(limit int64) {
	c.conn.SetReadLimit(limit)
}

func (c *connectionImpl) Subprotocol() string {
	return c.conn.Subprotocol()
}

func (c *connectionImpl) UnderlyingConn() net.Conn {
	return c.conn.UnderlyingConn()
}

func (c *connectionImpl) SetPingHandler(h func(appData string) error) {
	c.conn.SetPingHandler(h)
}

func (c *connectionImpl) SetPongHandler(h func(appData string) error) {
	c.conn.SetPongHandler(h)
}

func (c *connectionImpl) SetCloseHandler(h func(code int, text string) error) {
	c.conn.SetCloseHandler(h)
}

func (c *connectionImpl) CloseHandler() func(code int, text string) error {
	return c.conn.CloseHandler()
}

func (c *connectionImpl) PingHandler() func(appData string) error {
	return c.conn.PingHandler()
}

func (c *connectionImpl) PongHandler() func(appData string) error {
	return c.conn.PongHandler()
}
