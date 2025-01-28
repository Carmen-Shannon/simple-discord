package ws

import (
	"context"
	"net"

	"github.com/coder/websocket"
)

type Connection interface {
	Read(ctx context.Context) ([]byte, error)
	Write(ctx context.Context, data []byte, binary bool) error
	Close(statusCode int, reason string) error
	CloseNow() error
}

type WsConn struct {
	conn *websocket.Conn
}

func NewWebSocketConn(conn *websocket.Conn) *WsConn {
	return &WsConn{conn: conn}
}

func (ws *WsConn) Read(ctx context.Context) ([]byte, error) {
	_, data, err := ws.conn.Read(ctx)
	return data, err
}

func (ws *WsConn) Write(ctx context.Context, data []byte, binary bool) error {
	if binary {
		return ws.conn.Write(ctx, websocket.MessageBinary, data)
	}

	return ws.conn.Write(ctx, websocket.MessageText, data)
}

func (ws *WsConn) Close(statusCode int, reason string) error {
	return ws.conn.Close(websocket.StatusCode(statusCode), reason)
}

func (ws *WsConn) CloseNow() error {
	return ws.conn.CloseNow()
}

func (ws *WsConn) SetReadLimit(limit int64) {
	ws.conn.SetReadLimit(limit)
}

type UdpConn struct {
	conn    *net.UDPConn
	readBuf []byte
}

func NewUdpConn(conn *net.UDPConn) *UdpConn {
	return &UdpConn{
		conn:    conn,
		readBuf: make([]byte, 4096),
	}
}

func (u *UdpConn) Read(ctx context.Context) ([]byte, error) {
	n, _, err := u.conn.ReadFromUDP(u.readBuf)
	return u.readBuf[:n], err
}

func (u *UdpConn) Write(ctx context.Context, data []byte, binary bool) error {
	_, err := u.conn.Write(data)
	return err
}

func (u *UdpConn) Close(statusCode int, reason string) error {
	return u.conn.Close()
}

func (u *UdpConn) CloseNow() error {
	return u.conn.Close()
}
