/*
* Package tcpclient is a TCP client that provides a simple abstraction for reading and writing from TCP stream.me chat servers.
*
* Features:
*	- Retries happen automatically on disconnect.
*	- Thread safe
*
* Assumptions:
*
*   - messages arrive in a `<COMMAND> key="value" key2="value2"` ending in a '\n' character format
 */
package tcpclient

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"io"
	"log"
	"math/big"
	"net"
	"sync"
	"time"

	"github.com/StreamMeBots/pkg/commands"
	"github.com/StreamMeBots/pkg/parser"
)

// Client represents methods used to interface with the connected TCP server
type Client struct {
	// client host
	host string

	// comm channels
	stop  chan struct{}
	read  chan *readCommand
	write chan *writeMsg

	// state vars
	stats *Stats

	sync.RWMutex

	// used add a wait to the Read loop
	wg      sync.WaitGroup // used to pause reads when connection is lost
	waiting bool
}

// Stats represent some basic stats and state of a Client
type Stats struct {
	Online     bool          `json:"online"`
	Err        string        `json:"error,omitempty"`
	RetryCount int           `json:"retryCount"`
	Started    time.Time     `json:"started"`
	Uptime     time.Duration `json:"duration"`
}

type writeMsg struct {
	timeout time.Duration
	message string
	err     chan error
}

type readCommand struct {
	timeout time.Duration
	command chan *commands.Command
	err     chan error
}

// New is the constructor for Client. This function is required to create a Client.
func New(chatServerHost string) *Client {
	c := &Client{
		read:  make(chan *readCommand, 10),
		write: make(chan *writeMsg, 10),
		stop:  make(chan struct{}),
		host:  chatServerHost,
		stats: &Stats{},
	}

	go c.run()

	return c
}

// Close is used to gracefully shutdown the client
func (c *Client) Close() error {
	close(c.stop)
	return nil
}

// Write writes a message to the connected TCP server. A writeTimeout of 0 indicates no timeout
//
// NOTE: '\n' character is added to the message
func (c *Client) Write(msg string, writeTimeout time.Duration) error {
	c.wg.Wait()
	m := &writeMsg{
		timeout: writeTimeout,
		message: msg,
		err:     make(chan error),
	}
	c.write <- m
	err := <-m.err
	close(m.err)
	return err
}

// Stats returns stats about the tcp connection
func (c *Client) Stats() Stats {
	c.RLock()
	defer c.RUnlock()

	s := Stats{
		Online:     c.stats.Online,
		Started:    c.stats.Started,
		RetryCount: c.stats.RetryCount,
	}
	if !s.Started.IsZero() {
		s.Uptime = time.Since(c.stats.Started)
	}

	return s
}

// Read is used to read from the TCP server. This method should be called in a loop. A
// readTimeout of zero indicates no timeout.
func (c *Client) Read(readTimeout time.Duration) (*commands.Command, error) {
	c.wg.Wait()
	r := &readCommand{
		timeout: readTimeout,
		err:     make(chan error),
		command: make(chan *commands.Command),
	}
	c.read <- r
	select {
	case cmd := <-r.command:
		return cmd, nil
	case err := <-r.err:
		if err == io.EOF {
			c.wg.Add(1)
			c.waiting = true
		}
		return nil, err
	}
}

func (c *Client) run() error {
	var conn net.Conn
	var err error

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	// write loop
	go func() {
		for {
			select {
			case <-c.stop:
				return
			case w := <-c.write:
				if w.timeout > 0 {
					conn.SetWriteDeadline(time.Now().Add(w.timeout))
				} else {
					conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
				}
				r := bufio.NewWriter(conn)
				if _, err := r.WriteString(w.message + "\n"); err != nil {
					w.err <- err
					continue
				}
				if err := r.Flush(); err != nil {
					w.err <- err
					continue
				}
				w.err <- nil
			}
		}
	}()

	// read loop, also reconnects the tcp connection
	for {
		select {
		case <-c.stop:
			return nil
		default:
			conn, err = tls.Dial("tcp", c.host, &tls.Config{InsecureSkipVerify: true})
			if err != nil {
				c.connErr(err)
				continue
			}

			c.connected()

			stop := false
			for !stop {
				select {
				case <-c.stop:
					return nil
				case r := <-c.read:
					if r.timeout > 0 {
						conn.SetReadDeadline(time.Now().Add(r.timeout))
					} else {
						// no timeout
						conn.SetReadDeadline(time.Time{})
					}

					cmd, err := parser.Parse(conn)
					if err == io.EOF {
						stop = true
						c.connErr(err)
						r.err <- err
						break
					}
					if err != nil {
						r.err <- err
					} else {
						r.command <- cmd
					}
				}
			}
		}
	}
}

// connected sets the Client's state vars
func (c *Client) connected() {
	c.Lock()
	defer c.Unlock()

	c.stats.Online = true
	c.stats.Started = time.Now().UTC()
	c.stats.Err = ""
	if c.waiting {
		c.wg.Done()
	}
}

// connErr sets the Client's state vars
func (c *Client) connErr(err error) {
	c.Lock()
	defer c.Unlock()

	log.Printf("tcpclient: connection error %v. Trying again", err)

	c.stats.Online = false
	c.stats.Started = time.Time{}
	c.stats.Err = err.Error()
	c.stats.RetryCount++

	c.retrySleep()
}

// random sleep between 1 and 3 seconds
func (c *Client) retrySleep() {
	time.Sleep(randDuration(time.Second, time.Second*3))
}

func randDuration(min, max time.Duration) time.Duration {
	bg := big.NewInt(int64(max) - int64(min))

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return max / min
	}

	// add n to min to support the passed in range
	return time.Duration(n.Int64()) + time.Duration(min)
}
