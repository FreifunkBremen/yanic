// Importer for global RRD stats
package rrd

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bdlm/log"
)

var linePattern = regexp.MustCompile("^<!-- ....-..-.. ..:..:.. [A-Z]+ / (\\d+) --> <row><v>([^<]+)</v><v>([^<]+)</v></row>")

// Dataset a timestemp with values (node and clients)
type Dataset struct {
	Time    time.Time
	Nodes   float64
	Clients float64
}

// Read a rrdfile and return a chanel of datasets
func Read(rrdFile string) chan Dataset {
	out := make(chan Dataset)
	cmd := exec.Command("rrdtool", "dump", rrdFile)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Panicf("error on get stdout: %s", err)
	}
	if err := cmd.Start(); err != nil {
		log.Panicf("error on start rrdtool: %s", err)
	}

	r := bufio.NewReader(stdout)
	found := false

	go func() {
		for {
			// Read stdout by line
			line, _, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			str := strings.TrimSpace(string(line))

			// Search for the start of the daily datasets
			if !found {
				found = strings.Contains(str, "<!-- 86400 seconds -->")
				continue
			}
			if matches := linePattern.FindStringSubmatch(str); matches != nil && matches[2] != "NaN" && matches[3] != "NaN" {
				seconds, _ := strconv.Atoi(matches[1])
				nodes, _ := strconv.ParseFloat(matches[2], 64)
				clients, _ := strconv.ParseFloat(matches[3], 64)

				out <- Dataset{
					Time:    time.Unix(int64(seconds), 0),
					Nodes:   nodes,
					Clients: clients,
				}
			}
		}
		close(out)
	}()
	return out
}
