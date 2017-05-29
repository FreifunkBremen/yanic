package socket

type Message struct {
	Event string      `json:"event"`
	Body  interface{} `json:"body,omitempty"`
}

const (
	MessageEventInsertNode    = "insert_node"
	MessageEventInsertGlobals = "insert_globals"
	MessageEventPruneNodes    = "prune_nodes"
)
