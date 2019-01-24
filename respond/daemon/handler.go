package respondd

import (
	"encoding/json"
	"net"
	"reflect"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/respond"
)

func (d *Daemon) handler(socket *net.UDPConn) {
	socket.SetReadBuffer(respond.MaxDataGramSize)

	// Loop forever reading from the socket
	for {
		buf := make([]byte, respond.MaxDataGramSize)
		n, src, err := socket.ReadFromUDP(buf)
		if err != nil {
			log.Errorf("ReadFromUDP failed: %s", err)
		}
		raw := make([]byte, n)
		copy(raw, buf)

		get := string(raw)

		data := d.getData(src.Zone)

		log.WithFields(map[string]interface{}{
			"bytes": n,
			"data":  get,
			"src":   src.String(),
		}).Debug("recieve request")

		if get[:3] == "GET" {
			res, err := respond.NewRespone(data, src)
			if err != nil {
				log.Errorf("Decode failed: %s", err)
				continue
			}
			n, err = socket.WriteToUDP(res.Raw, res.Address)
			if err != nil {
				log.Errorf("WriteToUDP failed: %s", err)
				continue
			}
			log.WithFields(map[string]interface{}{
				"bytes": n,
				"dest":  res.Address.String(),
			}).Debug("send respond")
			continue
		}

		found := false

		t := reflect.TypeOf(data).Elem()
		v := reflect.ValueOf(data).Elem()

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fv := v.FieldByName(f.Name)
			if f.Tag.Get("json") == get {
				log.WithFields(map[string]interface{}{
					"param": get,
					"dest":  src.String(),
				}).Debug("found")
				raw, err = json.Marshal(fv.Interface())
				found = true
				break
			}
		}

		if !found {
			log.WithFields(map[string]interface{}{
				"param": get,
				"dest":  src.String(),
			}).Debug("not found")
			raw = []byte("ressource not found")
		}

		n, err = socket.WriteToUDP(raw, src)
		if err != nil {
			log.Errorf("WriteToUDP failed: %s", err)
			continue
		}
		log.WithFields(map[string]interface{}{
			"bytes": n,
			"dest":  src.String(),
		}).Debug("send respond")
	}
}
