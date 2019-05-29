package respondd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Vivena/babelweb2/state"
	"github.com/bdlm/log"
)

func (d *Daemon) babelConnect() {
	var conn net.Conn
	var err error

	for {
		log.Debug("Trying", d.Babel)
		for {
			conn, err = net.Dial("tcp6", d.Babel)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * 5)
			} else {
				break
			}
		}
		log.Info("Connected to ", d.Babel)
		closeConn := func() {
			conn.Close()
			log.Infof("Connection to %v closed\n", d.Babel)
		}

		fmt.Fprintf(conn, "monitor\n")
		s, err := state.NewBabelState(bufio.NewReader(conn), 0)
		if err == io.EOF {
			log.Warnf("Something wrong with %v:\n\tcouldn't get router id.\n", d.Babel)
		} else if err != nil {
			// Don't you even dare to reconnect to this unholy node!
			log.Warnf("Oh, boy! %v is doomed:\n\t%v.\t", d.Babel, err)
			closeConn()
			return
		} else {
			d.babelData = s
			err := s.ListenHistory()
			if err != nil {
				log.Warnf("Babel listen stopped: %s", err)
			}
			d.babelData = nil
		}
		closeConn()
	}
}
