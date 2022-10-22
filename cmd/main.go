package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

//go:embed templates/index.html
var indexTmpl embed.FS

func main() {
	http.HandleFunc("/", index)
	http.Handle("/ws", websocket.Handler(msgHandler))

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

func msgHandler(ws *websocket.Conn) {
	defer ws.Close()

	// 初回のメッセージを送信
	err := websocket.Message.Send(ws, "こんにちは！ :)")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		// メッセージを受信する
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		log.Println(msg)
		if err != nil {
			log.Fatalln(err)
		}

		// メッセージを返信する
		err := websocket.Message.Send(ws, fmt.Sprintf(`%q というメッセージを受け取りました。`, msg))
		if err != nil {
			log.Fatalln(err)
		}
	}
}
