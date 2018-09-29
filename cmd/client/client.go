package main

import (
	"log"
	"github.com/DGHeroin/gltun"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client := gltun.NewClient()

	client.HandleDataFunc(onData)
	client.HandleClosedFunc(onClosed)

	if session, err := client.Connect("127.0.0.1:9999"); err == nil {
		go func() {
			session.Send([]byte("test"))
		}()
		if err := client.Wait(); err != nil {
			log.Println(err)
		}
	}
}

func onConnected(session gltun.Session) {
	log.Println("on connected:", session)
}
func onClosed(session gltun.Session) {
	log.Println("on closed:", session)
}


func onData(session gltun.Session, data []byte) {
	log.Println(string(data))

	if err := session.Send([]byte("client say world")); err != nil {
		log.Println(err)
	}
}