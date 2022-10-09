package internal

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// CloseMessageCode identifier the message id for a close request
	CloseMessageCode = 1000
	// DefaultTimeout timeout in seconds
	DefaultTimeout = 30
)

type Option func(ws *WebSocket) error

type WebSocket struct {
	conn    *websocket.Conn
	timeout time.Duration

	respChan chan *RPCRawResponse
	close    chan int
}

// NewWebsocket creates a new websocket connection
// If timeout is not specified, it will use the default of 30s
func NewWebsocket(url string, timeout ...*DbTimeoutConfig) (*WebSocket, error) {
	ws, err := NewWebsocketWithOptions(url, Timeout(DefaultTimeout))
	if err != nil {
		return nil, err
	}
	if len(timeout) == 0 || timeout[0] == nil {
		return ws, nil
	}

	ws.timeout = timeout[0].Timeout

	return ws, nil
}

func NewWebsocketWithOptions(url string, options ...Option) (*WebSocket, error) {
	dialer := websocket.DefaultDialer
	dialer.EnableCompression = true

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	ws := &WebSocket{
		conn:     conn,
		respChan: make(chan *RPCRawResponse),
		close:    make(chan int),
	}

	for _, option := range options {
		if err := option(ws); err != nil {
			return nil, err
		}
	}

	ws.initialize()
	return ws, nil
}

func Timeout(timeout float64) Option {
	return func(ws *WebSocket) error {
		ws.timeout = time.Duration(timeout) * time.Second
		return nil
	}
}

func (ws *WebSocket) Close() error {
	defer func() {
		close(ws.close)
	}()

	return ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(CloseMessageCode, ""))
}

func (ws *WebSocket) Send(id, method string, params []interface{}) (*RPCRawResponse, error) {
	request := &RPCRequest{
		ID:     id,
		Method: method,
		Params: params,
	}

	if err := ws.write(request); err != nil {
		return nil, err
	}

	tick := time.NewTicker(ws.timeout)

	for {
		select {
		case <-tick.C:
			return nil, errors.New("timeout")
		case res := <-ws.respChan:
			if res.id != id {
				continue
			}
			if res.HasInternalError() {
				return nil, res.Error()
			}

			return res, nil
		}
	}
}

func (ws *WebSocket) read() (*RPCRawResponse, error) {
	_, data, err := ws.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return CreateRPCRawResponse(data), nil
}

func (ws *WebSocket) write(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

func (ws *WebSocket) initialize() {
	go func() {
		for {
			select {
			case <-ws.close:
				return
			default:
				res, err := ws.read()

				if err != nil {
					log.Println("error reading from websocket:", err)
					continue
				}

				ws.respChan <- res
			}
		}
	}()
}
