package tcp

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type handle func(string, *net.Conn)

type Server struct {
	listener   net.Listener
	serverAddr string
	running    bool
	FnConn     handle
	startTime  time.Time
}

func New(host string, handle func(string, *net.Conn)) (*Server, error) {
	s := &Server{
		serverAddr: host,
		FnConn:     handle,
	}

	return s, s.Start()
}

func (s *Server) Start() error {
	var err error
	s.startTime = time.Now()
	s.running = true
	s.listener, err = net.Listen("tcp", s.serverAddr)
	if err != nil {
		return err
	}
	defer s.listener.Close()

	go s.FnConn("server", nil)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()

		go s.handleConnection(conn)
	}
	return nil
}

func (s *Server) Stop() error {
	if err := s.listener.Close(); err != nil {
		return fmt.Errorf("failed to stop server: %v", err)
	}
	s.running = false
	return nil
}

func (s *Server) Refresh() error {
	if err := s.Stop(); err != nil {
		return err
	}
	return s.Start()
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		message = message[:len(message)-1]

		switch message {
		case "-stop":
			fmt.Fprintln(conn, "Server stopping...")
			s.Stop()
		case "-refresh":
			if err := s.Refresh(); err == nil {
				fmt.Fprintln(conn, "Server refreshing...")
			}
		case "-status":
			fmt.Fprintln(conn, "Server running:", s.running)
		case "-uptime":
			if s.running {
				fmt.Fprintln(conn, "Server uptime:", time.Since(s.startTime))
			} else {
				fmt.Fprintln(conn, "Server is not running")
			}
		default:
			s.FnConn(message, &conn)
		}
	}

}
