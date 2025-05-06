package readers

import (
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
)

type udpReader struct {
	logger *zap.Logger
	conn   *net.UDPConn
}

type udpoption struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

type udpOption func(*udpoption)

func WithUDPConnection(conn *net.UDPConn) udpOption {
	return func(opt *udpoption) {
		opt.conn = conn
	}
}

func WithUDPAddr(addr *net.UDPAddr) udpOption {
	return func(opt *udpoption) {
		opt.addr = addr
	}
}

func NewUDPReader(logger *zap.Logger, opts ...udpOption) (*udpReader, error) {
	var errs error
	// Validate the inputs and accumulate errors.
	if logger == nil {
		errs = errors.Join(errs, errors.New("logger cannot be nil"))
	}

	opt := &udpoption{}
	for _, o := range opts {
		o(opt)
	}

	if opt.conn == nil && opt.addr == nil {
		errs = errors.Join(errs, errors.New("bad udp config"))
	}

	if errs != nil {
		return nil, errs
	}

	var conn *net.UDPConn = opt.conn
	var err error

	if opt.addr != nil {
		logger.Info("Address provided, creating new UDP Connection", zap.Any("Address", opt.addr))
		conn, err = net.ListenUDP("udp", opt.addr)
		if err != nil {
			return nil, err
		}
	}

	return &udpReader{
		logger: logger,
		conn:   conn,
	}, nil
}

// Read implements the Read method that adheres to io.Reader.
func (r *udpReader) Read(ctx context.Context, b []byte) (int, error) {
	n, _, err := r.conn.ReadFromUDP(b)
	r.logger.Debug("Received message from UDP connection")
	if err != nil {
		r.logger.Error("error reading connection", zap.Error(err))
		return n, err
	}

	// TODO: record metrics here if you have time

	return n, nil
}
