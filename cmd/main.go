package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var clients = make(map[*websocket.Conn]bool)

// メッセージブロードキャストチャネル
var broadcast = make(chan string)

//go:embed templates/index.html
var indexTmpl embed.FS

func main() {
	http.HandleFunc("/", index)
	http.Handle("/ws", websocket.Handler(handleConnection))
	go handleMessages()

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(indexTmpl, "templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		log.Fatal(err)
	}
}

func handleConnection(ws *websocket.Conn) {
	defer ws.Close()

	// 初回のメッセージを送信
	err := websocket.Message.Send(ws, "こんにちは！")
	if err != nil {
		log.Fatalln(err)
	}

	clients[ws] = true
	for {
		// メッセージを受信する
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Fatalln(err)
		}
		// 受け取ったメッセージをbroadcastチャネルに送る
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// broadcastチャネルからメッセージを受け取る
		msg := <-broadcast
		// 接続中の全クライアントにメッセージを送る
		for client := range clients {
			// メッセージを返信する
			err := websocket.Message.Send(client, fmt.Sprintf(`%q というメッセージを受け取りました。`, msg))
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
