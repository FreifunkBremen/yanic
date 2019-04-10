package respondd

import (
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/yanic/data"
)

type Batman struct {
	Bridge     string
	Interfaces []string
}

func NewBatman(iface string) *Batman {
	out, err := exec.Command("batctl", "-m", iface, "if").Output()
	if err != nil {
		log.WithField("iface", iface).Error("not able to run batctl")
		return nil
	}
	b := &Batman{Bridge: iface}
	for _, line := range strings.Split(string(out), "\n") {
		i := strings.Split(line, ":")[0]
		if i != "" {
			b.Interfaces = append(b.Interfaces, i)
		}
	}
	return b
}
func (b *Batman) Address(iface string) string {
	i, err := net.InterfaceByName(iface)
	if err != nil {
		return ""
	}
	return i.HardwareAddr.String()
}
func (b *Batman) Neighbours() map[string]data.BatadvNeighbours {
	out, err := exec.Command("batctl", "-m", b.Bridge, "o").Output()
	if err != nil {
		log.WithField("iface", b.Bridge).Error("not able to run batctl")
		return nil
	}

	lines := strings.Split(string(out), "\n")
	neighbours := make(map[string]data.BatadvNeighbours)

	re := regexp.MustCompile(`([0-9a-f:]+)\s+(\d+\.\d+)s\s+\((\d+)\)\s+([0-9a-f:]+)\s+\[\s*([a-z0-9-]+)\]`)

	for _, i := range b.Interfaces {
		mac := b.Address(i)
		neighbour := data.BatadvNeighbours{
			Neighbours: make(map[string]data.BatmanLink),
		}

		for _, line := range lines {
			fields := re.FindStringSubmatch(line)
			if len(fields) != 6 {
				continue
			}
			if i == fields[5] && fields[1] == fields[4] {
				lastseen, err := strconv.ParseFloat(fields[2], 64)
				if err != nil {
					log.WithField("value", fields[2]).Warnf("unable to parse lastseen: %s", err)
					continue
				}

				tq, err := strconv.Atoi(fields[3])
				if err != nil {
					log.WithField("value", fields[3]).Warnf("unable to parse tq: %s", err)
					continue
				}

				nMAC := fields[1]

				neighbour.Neighbours[nMAC] = data.BatmanLink{
					Lastseen: lastseen,
					Tq:       tq,
				}
			}
		}
		neighbours[mac] = neighbour
	}
	return neighbours
}
