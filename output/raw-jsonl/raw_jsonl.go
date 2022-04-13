package jsonlines

import (
	"yanic/data"
	"yanic/lib/jsontime"
	"yanic/runtime"
)

// RawNode struct
type RawNode struct {
	Firstseen    jsontime.Time          `json:"firstseen"`
	Lastseen     jsontime.Time          `json:"lastseen"`
	Online       bool                   `json:"online"`
	Statistics   *data.Statistics       `json:"statistics"`
	Nodeinfo     *data.Nodeinfo         `json:"nodeinfo"`
	Neighbours   *data.Neighbours       `json:"neighbours"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type FileInfo struct {
	Version   int           `json:"version"`
	Timestamp jsontime.Time `json:"updated_at"` // Timestamp of the generation
	Format    string        `json:"format"`
}

func transform(nodes *runtime.Nodes) []interface{} {
	var result []interface{}

	result = append(result, FileInfo{
		Version:   1,
		Timestamp: jsontime.Now(),
		Format:    "raw-nodes-jsonl",
	})

	for _, nodeOrigin := range nodes.List {
		if nodeOrigin != nil {
			node := &RawNode{
				Firstseen:    nodeOrigin.Firstseen,
				Lastseen:     nodeOrigin.Lastseen,
				Online:       nodeOrigin.Online,
				Statistics:   nodeOrigin.Statistics,
				Nodeinfo:     nodeOrigin.Nodeinfo,
				Neighbours:   nodeOrigin.Neighbours,
				CustomFields: nodeOrigin.CustomFields,
			}
			result = append(result, node)
		}
	}
	return result
}
