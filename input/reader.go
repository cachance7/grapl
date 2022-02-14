package input

import (
	"bufio"
	"bytes"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// Message contains the input payload
type Message struct {
	Data []byte
}

// Reader can Start and Stop reading
type Reader interface {
	Start() error
	Stop() error
	Read() (Message, error)
}

// StdinReader reads from stdin
type StdinReader struct {
	buf  chan []byte
	out  chan Message
	done bool
}

func (reader StdinReader) Read() (Message, error) {
	msg, ok := <-reader.out
	if !ok {
		log.Println("couldn't read from buf")
	}
	log.Println("got a message for Read()")
	return msg, nil
}

// Start begins the reader loop
func (reader StdinReader) Start() {
	go reader.readLoop()
	go reader.writeLoop()
}

// Stop terminates the reader
func (reader *StdinReader) Stop() {
	reader.done = true
}

// NewStdinReader creates a new StdinReader ready to be Start()ed
func NewStdinReader() StdinReader {
	return StdinReader{buf: make(chan []byte, 100), out: make(chan Message, 100)}
}

func (reader StdinReader) writeLoop() {
	var buf bytes.Buffer
	for {
		select {
		case msg := <-reader.buf:
			log.Printf("copying %d bytes to buffer", len(msg))
			buf.Write(msg)
		case <-time.After(300 * time.Millisecond):
			if buf.Len() > 0 {
				log.Printf("writing Message to out buffer")
				outbuf := make([]byte, buf.Len())
				copy(outbuf, buf.Bytes())
				buf.Reset()
				reader.out <- Message{outbuf}
			}
		}
	}
}

// readLoop continually reads from stdin and pushes chunks of bytes into the buf channel
func (reader StdinReader) readLoop() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		text := scanner.Bytes()

		if len(text) != 0 {
			log.Printf("writing %d bytes from stdin\n", len(text))
			reader.buf <- text
		}

		if reader.done {
			log.Println("terminating readLoop")
			break
		}
	}
}
