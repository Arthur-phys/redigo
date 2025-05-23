//go:build integration
// +build integration

package server

import (
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Arthur-phys/redigo/pkg/core/cache"
	"github.com/Arthur-phys/redigo/pkg/core/respparser"
)

func TestIntegration_WorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_A_Single_One(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 10240),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection

	go func() {
		newWorker.handleConnection(&genericConn)
	}()
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n"))
	response := make([]byte, 1024)
	n, _ := newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected message received! %v - %v", string(response), response)
	}
}

func TestIntegration_WorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_Multiple_In_A_Single_Segment(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 10240),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	response := make([]byte, 1024)
	n, _ := newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n$7\r\ncrayoli\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}
}

func TestIntegration_WorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_Multiple_Commands_In_Multiple_Segments(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 10240),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r"))
	response := make([]byte, 1024)
	n, _ := newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "\n$3\r\nGET\r\n$1\r\nB\r\n"))
	response = make([]byte, 1024)
	n, _ = newConnection.readAsClient(response)
	if string(response[:n]) != "$7\r\ncrayoli\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}
}

func TestIntegration_WorkerhandleConnection_Should_Return_Message_To_Client_When_Sent_Multiple_Commands_In_Even_More_Segments(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 10240),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r"))
	response := make([]byte, 1024)
	n, _ := newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "\n$3\r\nGET\r\n$1\r\n"))
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "B\r\n*4\r\n$5\r\nLPUSH\r\n$4\r\nCats\r\n$4\r\nNiji\r\n$7\r\nBigotes\r\n"))

	response = make([]byte, 1024)
	n, _ = newConnection.readAsClient(response)
	if string(response[:n]) != "$7\r\ncrayoli\r\n_\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}
}

func TestIntegration_WorkerhandleConnection_Should_Return_Message_To_Client_When_Partitioned_In_Different_Ways(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 10240),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n"))
	response := make([]byte, 1024)
	n, _ := newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}

	// Partitioned multiple commands
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\n"))
	// A command partitioned in the array declaration
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "B\r\n*4\r"))
	response = make([]byte, 1024)
	n, _ = newConnection.readAsClient(response)
	if string(response[:n]) != "$7\r\ncrayoli\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}
	// Then in the raw string declaration
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "\n$5\r"))
	// Then in the raw string content
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "\nLPU"))
	// Then the rest of the command
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "SH\r\n$4\r\nCats\r\n$4\r\nNiji\r\n$7\r\nBigotes\r\n"))

	response = make([]byte, 1024)
	n, _ = newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}
}

func TestIntegration_WorkerhandleConnection_Should_Return_Error_To_Client_When_Sent_Multiple_Commands_With_One_Wrong(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 10240),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*1\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	response := make([]byte, 1024)
	newConnection.readAsClient(response)
	if response[0] != '-' {
		t.Errorf("Unexpected response! %v", string(response))
	}
}

func TestIntegration_WorkerhandleConnection_Should_Return_Error_To_Client_When_Exceeding_Command_Size(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 50),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	response := make([]byte, 1024)
	n, err := newConnection.readAsClient(response)
	if string(response[:n]) != "-Call exceeded size allowed\r\n" {
		t.Errorf("Unexpected message received! %v - %v", string(response), err)
	}
}

func TestIntegration_WorkerhandleConnection_Should_Return_Error_To_Client_When_Exceeding_Command_Size_In_Subsecuent_Command(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 36),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n*2"))
	response := make([]byte, 1024)
	n, _ := newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}

	response = make([]byte, 1024)
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	n, err := newConnection.readAsClient(response)
	if string(response[:n]) != "-Call exceeded size allowed\r\n" {
		t.Errorf("Unexpected message received! %v - %v", string(response), err)
	}

}

func TestIntegration_WorkerhandleConnection_Should_Return_Message_To_Client_When_Exceeding_Command_Size_But_Sending_Commands_Sepparated(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cacheStore := cache.New()
	channel := make(chan net.Conn)
	newWorker := worker{
		cacheStore:     cacheStore,
		connections:    channel,
		timeout:        1,
		notifications:  make(chan struct{}, 1),
		id:             1,
		parser:         respparser.New(nil, 36),
		shutdownWaiter: &sync.WaitGroup{},
	}

	var genericConn net.Conn
	newConnection := newMockConnection()
	defer newConnection.Close()
	genericConn = &newConnection
	go func() {
		newWorker.handleConnection(&genericConn)
	}()

	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*3\r\n$3\r\nSET\r\n$1\r\nB\r\n$7\r\ncrayoli\r\n"))
	response := make([]byte, 1024)
	n, _ := newConnection.readAsClient(response)
	if string(response[:n]) != "_\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}

	response = make([]byte, 1024)
	newConnection.writeAsClient(fmt.Appendf([]byte{}, "*2\r\n$3\r\nGET\r\n$1\r\nB\r\n"))
	n, _ = newConnection.readAsClient(response)
	if string(response[:n]) != "$7\r\ncrayoli\r\n" {
		t.Errorf("Unexpected message received! %v", string(response))
	}

}

type mockConnection struct {
	requestArray        []byte
	responseArray       []byte
	currentBytesRead    int
	closed              bool
	newData             bool
	requestMutex        *sync.Mutex
	requestConditional  *sync.Cond
	responseMutex       *sync.Mutex
	responseConditional *sync.Cond
}

func newMockConnection() mockConnection {
	mc := mockConnection{
		requestArray:     []byte{},
		responseArray:    []byte{},
		currentBytesRead: 0,
		closed:           false,
		newData:          true,
		requestMutex:     &sync.Mutex{},
		responseMutex:    &sync.Mutex{},
	}
	mc.requestConditional = sync.NewCond(mc.requestMutex)
	mc.responseConditional = sync.NewCond(mc.responseMutex)
	return mc
}

func (mc *mockConnection) Read(b []byte) (int, error) {
	mc.requestMutex.Lock()
	defer mc.requestMutex.Unlock()

	for !mc.closed && len(mc.requestArray) == 0 {
		mc.requestConditional.Wait()
	}
	if mc.closed {
		return 0, io.EOF
	}
	n := int(math.Min(float64(len(b)), float64(len(mc.requestArray))))
	copy(b, mc.requestArray[:n])
	mc.requestArray = mc.requestArray[n:]
	return n, nil
}

func (mc *mockConnection) writeAsClient(b []byte) (int, error) {
	mc.requestMutex.Lock()
	defer mc.requestMutex.Unlock()
	defer mc.requestConditional.Broadcast()
	if mc.closed {
		return 0, io.EOF
	}
	mc.requestArray = append(mc.requestArray, b...)
	return len(b), nil
}

func (mc *mockConnection) Write(b []byte) (int, error) {
	mc.responseMutex.Lock()
	defer mc.responseMutex.Unlock()
	defer mc.responseConditional.Broadcast()
	if mc.closed {
		return 0, io.EOF
	}
	mc.responseArray = append(mc.responseArray, b...)
	return len(b), nil
}

func (mc *mockConnection) readAsClient(b []byte) (int, error) {
	mc.responseMutex.Lock()
	defer mc.responseMutex.Unlock()

	for !mc.closed && len(mc.responseArray) == 0 {
		mc.responseConditional.Wait()
	}
	n := int(math.Min(float64(len(b)), float64(len(mc.responseArray))))
	copy(b, mc.responseArray[:n])
	mc.responseArray = mc.responseArray[n:]
	return n, nil
}

func (mc *mockConnection) Close() error {
	mc.requestMutex.Lock()
	mc.responseMutex.Lock()
	defer mc.requestMutex.Unlock()
	defer mc.responseMutex.Unlock()
	mc.closed = true
	mc.requestConditional.Broadcast()
	mc.responseConditional.Broadcast()
	return nil
}

func (mc *mockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (mc *mockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}
func (mc *mockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (mc *mockConnection) LocalAddr() net.Addr {
	return mockAddr{}
}

func (mc *mockConnection) RemoteAddr() net.Addr {
	return mockAddr{}
}

type mockAddr struct{}

func (ma mockAddr) Network() string {
	return "no network"
}

func (ma mockAddr) String() string {
	return "test"
}
