package tcp

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	xnet "github.com/go-pantheon/fabrica-net"
	"github.com/go-pantheon/fabrica-net/xcontext"
	"github.com/go-pantheon/fabrica-net/internal/bufreader"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var ErrTimeout = errors.New("i/o timeout")

type Option func(c *Client)

func Bind(bind string) Option {
	return func(s *Client) {
		s.bind = bind
	}
}

type Client struct {
	xsync.Stoppable

	Id   int64
	bind string

	conn   *net.TCPConn
	reader *bufreader.Reader

	receivePackChan chan []byte
}

func NewClient(id int64, opts ...Option) *Client {
	c := &Client{
		Stoppable: xsync.NewStopper(time.Second * 10),
		Id:        id,
	}

	for _, o := range opts {
		o(c)
	}
	c.receivePackChan = make(chan []byte, 1024)

	return c
}

func (c *Client) Start(ctx context.Context) (err error) {
	addr, err := net.ResolveTCPAddr("tcp", c.bind)
	if err != nil {
		err = errors.Wrapf(err, "resolved failed. cli=%d addr=%s", c.Id, c.bind)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		err = errors.Wrapf(err, "connect failed. cli=%d addr=%s", c.Id, c.bind)
		return
	}

	xcontext.SetDeadlineWithContext(ctx, conn, fmt.Sprintf("client=%d", c.Id))

	c.reader = bufreader.NewReader(conn, 4096)
	c.conn = conn

	xsync.GoSafe(fmt.Sprintf("tcp.client.id=%d", c.Id), func() error {
		return c.receive(ctx)
	})
	return
}

func (c *Client) receive(ctx context.Context) error {
	defer c.stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-c.StopTriggered():
			c.stop()
			return xsync.GroupStopping
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		return xsync.RunSafe(func() error {
			return c.readPackLoop()
		})
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (c *Client) stop() {
	c.DoStop(func() {
		if c.reader != nil {
			if err := c.reader.Close(); err != nil {
				log.Errorf("[tcp.Client] cli=%d bufreader close failed", c.Id)
			}
		}
		close(c.receivePackChan)
		log.Debugf("[tcp.Client] cli=%d is closed", c.Id)
	})
}

func (c *Client) readPackLoop() error {
	for {
		pack, err := c.read()
		if err != nil {
			return err
		}
		c.receivePackChan <- pack
	}
}

func (c *Client) read() (buf []byte, err error) {
	lb, err := c.reader.ReadFull(xnet.PackLenSize)
	if err != nil {
		err = errors.Wrapf(err, "read pack len failed. lenSize=%d", xnet.PackLenSize)
		return
	}

	var packLen int32
	err = binary.Read(bytes.NewReader(lb), binary.BigEndian, &packLen)
	if err != nil {
		return
	}
	buf, err = c.reader.ReadFull(int(packLen))
	if err != nil {
		err = errors.Wrapf(err, "read pack body failed. len=%d", packLen)
	}
	return
}

func (c *Client) write(pack []byte) (err error) {
	var buf bytes.Buffer

	err = binary.Write(&buf, binary.BigEndian, uint32(len(pack)))
	if err != nil {
		err = errors.Wrapf(err, "write pack len failed. len=%d", len(pack))
		return
	}

	_, err = buf.Write(pack)
	if err != nil {
		err = errors.Wrapf(err, "write pack body failed. len=%d", len(pack))
		return
	}

	_, err = c.conn.Write(buf.Bytes())
	if err != nil {
		err = errors.Wrapf(err, "write pack failed. len=%d", len(pack))
		return
	}
	return
}

func (c *Client) Send(pack []byte) (err error) {
	return c.write(pack)
}

func (c *Client) Receive() <-chan []byte {
	return c.receivePackChan
}
