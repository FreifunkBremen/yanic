package respondd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Vivena/babelweb2/parser"
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
		defer closeConn()
		fmt.Fprintf(conn, "monitor\n")
		r := bufio.NewReader(conn)
		s := parser.NewScanner(r)
		desc := parser.NewBabelDesc()
		err = desc.Fill(s)
		if err == io.EOF {
			log.Warnf("Something wrong with %v:\n\tcouldn't get router id.\n", d.Babel)
		} else if err != nil {
			// Don't you even dare to reconnect to this unholy node!
			log.Warnf("Oh, boy! %v is doomed:\n\t%v.\t", d.Babel, err)
			return
		} else {
			d.babelData = desc
			err := d.babelDescListen(s)
			if err != nil {
				log.Warnf("Babel listen stopped: %s", err)
			}
			d.babelData = nil
		}
		closeConn()
	}
}

func (d *Daemon) babelDescListen(s *parser.Scanner) error {
	for {
		upd, err := d.babelData.ParseAction(s)
		if err != nil && err != io.EOF && err.Error() != "EOL" {
			return err
		}
		if err == io.EOF {
			break
		}
		//TODO maybe keep upd.action != none
		if !(d.babelData.CheckUpdate(upd)) {
			continue
		}
		err = d.babelData.Update(upd)
		if err != nil {
			return err
		}
	}
	return nil
}
