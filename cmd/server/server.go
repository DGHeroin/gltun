package main

import (
	"github.com/DGHeroin/gltun"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := gltun.NewServer()
	server.HandleConnectedFunc(onConnected)
	server.HandleClosedFunc(onClosed)
	server.HandleDataFunc(onData)
	server.ListenAndServe(":9999")
}

func onConnected(session gltun.Session) {
	log.Println("on connected:", session)
}
func onClosed(session gltun.Session) {
	log.Println("on closed:", session)
}

func onData(session gltun.Session, data []byte) {
	log.Println(string(data))
	session.Send(data)
}