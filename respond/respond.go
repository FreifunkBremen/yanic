package respond

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"

	"github.com/bdlm/log"
	"github.com/tidwall/gjson"

	"github.com/FreifunkBremen/yanic/data"
)

const (
	// default multicast group used by announced
	MulticastAddressDefault = "ff05:0:0:0:0:0:2:1001"

	// default udp port used by announced
	PortDefault = 1001

	// maximum receivable size
	MaxDataGramSize = 8192
)

// Response of the respond request
type Response struct {
	Address *net.UDPAddr
	Raw     []byte
}

func NewRespone(res *data.ResponseData, addr *net.UDPAddr) (*Response, error) {
	buf := new(bytes.Buffer)
	flater, err := flate.NewWriter(buf, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	defer flater.Close()

	if err = json.NewEncoder(flater).Encode(res); err != nil {
		return nil, err
	}

	err = flater.Flush()

	return &Response{
		Raw:     buf.Bytes(),
		Address: addr,
	}, err
}

func (res *Response) parse(customFields []CustomFieldConfig) (*data.ResponseData, error) {
	// Deflate
	deflater := flate.NewReader(bytes.NewReader(res.Raw))
	defer deflater.Close()

	jsonData, err := ioutil.ReadAll(deflater)
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	// Unmarshal
	rdata := &data.ResponseData{}
	err = json.Unmarshal(jsonData, rdata)

	rdata.CustomFields = make(map[string]interface{})
	if !gjson.Valid(string(jsonData)) {
		log.WithField("jsonData", jsonData).Info("JSON data is invalid")
	} else {
		jsonParsed := gjson.Parse(string(jsonData))
		for _, customField := range customFields {
			field := jsonParsed.Get(customField.Path)
			if field.Exists() {
				rdata.CustomFields[customField.Name] = field.String()
			}
		}
	}

	return rdata, err
}
