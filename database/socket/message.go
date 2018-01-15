package socket

type Message struct {
	Event string      `json:"event"`
	Body  interface{} `json:"body,omitempty"`
	Site  string      `json:"site"`
}

const (
	MessageEventInsertNode    = "insert_node"
	MessageEventInsertGlobals = "insert_globals"
	MessageEventInsertLink    = "insert_link"
	MessageEventPruneNodes    = "prune_nodes"
)
