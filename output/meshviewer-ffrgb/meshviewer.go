package meshviewerFFRGB

import (
	"fmt"
	"log"
	"strings"

	"github.com/FreifunkBremen/yanic/lib/jsontime"
	"github.com/FreifunkBremen/yanic/runtime"
)

func transform(nodes *runtime.Nodes) *Meshviewer {

	meshviewer := &Meshviewer{
		Timestamp: jsontime.Now(),
		Nodes:     make([]*Node, 0),
		Links:     make([]*Link, 0),
	}

	links := make(map[string]*Link)
	typeList := make(map[string]string)

	nodes.RLock()
	defer nodes.RUnlock()

	for _, nodeOrigin := range nodes.List {
		node := NewNode(nodes, nodeOrigin)
		meshviewer.Nodes = append(meshviewer.Nodes, node)

		if !nodeOrigin.Online {
			continue
		}

		if nodeinfo := nodeOrigin.Nodeinfo; nodeinfo != nil {
			if meshes := nodeinfo.Network.Mesh; meshes != nil {
				for _, mesh := range meshes {
					for _, mac := range mesh.Interfaces.Wireless {
						typeList[mac] = "wifi"
					}
					for _, mac := range mesh.Interfaces.Tunnel {
						typeList[mac] = "vpn"
					}
				}
			}
		}

		for _, linkOrigin := range nodes.NodeLinks(nodeOrigin) {
			var key string
			// keep source and target in the same order
			switchSourceTarget := strings.Compare(linkOrigin.SourceMAC, linkOrigin.TargetMAC) > 0
			if switchSourceTarget {
				key = fmt.Sprintf("%s-%s", linkOrigin.SourceMAC, linkOrigin.TargetMAC)
			} else {
				key = fmt.Sprintf("%s-%s", linkOrigin.TargetMAC, linkOrigin.SourceMAC)
			}
			if link := links[key]; link != nil {
				if switchSourceTarget {
					link.TargetTQ = float32(linkOrigin.TQ) / 255.0
					if link.Type == "other" {
						link.Type = typeList[linkOrigin.TargetMAC]
					} else if link.Type != typeList[linkOrigin.TargetMAC] {
						log.Printf("different linktypes %s:%s current: %s source: %s target: %s", linkOrigin.SourceMAC, linkOrigin.TargetMAC, link.Type, typeList[linkOrigin.SourceMAC], typeList[linkOrigin.TargetMAC])
					}
				} else {
					link.SourceTQ = float32(linkOrigin.TQ) / 255.0
					if link.Type == "other" {
						link.Type = typeList[linkOrigin.SourceMAC]
					} else if link.Type != typeList[linkOrigin.SourceMAC] {
						log.Printf("different linktypes %s:%s current: %s source: %s target: %s", linkOrigin.SourceMAC, linkOrigin.TargetMAC, link.Type, typeList[linkOrigin.SourceMAC], typeList[linkOrigin.TargetMAC])
					}
				}
				if link.Type == "" {
					link.Type = "other"
				}
				continue
			}
			tq := float32(linkOrigin.TQ) / 255.0
			link := &Link{
				Type:      typeList[linkOrigin.SourceMAC],
				Source:    linkOrigin.SourceID,
				SourceMAC: linkOrigin.SourceMAC,
				Target:    linkOrigin.TargetID,
				TargetMAC: linkOrigin.TargetMAC,
				SourceTQ:  tq,
				TargetTQ:  tq,
			}
			if switchSourceTarget {
				link.Type = typeList[linkOrigin.TargetMAC]
				link.Source = linkOrigin.TargetID
				link.SourceMAC = linkOrigin.TargetMAC
				link.Target = linkOrigin.SourceID
				link.TargetMAC = linkOrigin.SourceMAC
			}
			if link.Type == "" {
				link.Type = "other"
			}
			links[key] = link
			meshviewer.Links = append(meshviewer.Links, link)
		}
	}

	return meshviewer
}
