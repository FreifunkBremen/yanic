package main

import (
	"log"
	"net/http"
)

func main(){
	node := NewNodeServer("/nodes")
	go node.Listen()
	
	annouced := NewAnnouced(node)
	go annouced.Run()

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
