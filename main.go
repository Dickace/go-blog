package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ConnectUser struct {
	Websocket *websocket.Conn
	ClientIP string
	UserData User
}

type User struct {
	Username string
	Messages []Message
}

type requestUser struct {
	username string
}

type Message struct {
	Types string `json:"type"`
	Text  string `json:"text"`
}

func newConnectUser(ws *websocket.Conn, clientIP string,r *http.Request ) *ConnectUser  {
	var requestBody requestUser
	err:= json.NewDecoder(r.Body).Decode(&requestBody); if err!= nil {
		fmt.Print("can't decode Json")
	}
	return &ConnectUser{
		Websocket: ws,
		ClientIP: clientIP,
		UserData: User{
			Username: requestBody.username,
			Messages: []Message{},
		},
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request)  {
	tmpl, _ := template.ParseFiles("templates/index.html")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var users = []ConnectUser{}

func RemoveIndex(s []ConnectUser, index int) []ConnectUser {
	return append(s[:index], s[index+1:]...)
}



func WebsocketHandler(w http.ResponseWriter, r *http.Request)  {
	ws, _ := upgrader.Upgrade(w, r, nil)

	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()

	log.Println("Client connected:", ws.RemoteAddr().String())
	var socketClient *ConnectUser = newConnectUser(ws, ws.RemoteAddr().String(),r)
	users = append(users, *socketClient)
	log.Println("Number client connected ...", len(users))

	for {
		messageType, message, err := ws.ReadMessage()
		if  err != nil {
			log.Println("Ws disconnect waiting", err.Error())
			for i, v := range users {
				if v.ClientIP == socketClient.ClientIP {
					users = RemoveIndex(users,i)
				}
			}
			log.Println("Number of client still connected ...", len(users))
			return
		}
		log.Print(messageType)
		log.Print(message)
		var msg *Message
		err = json.Unmarshal(message,&msg); if err!=nil{
			fmt.Print("Decode error")
		}
		switch msg.Types {
		case "userList":
			break
		case "message":
			break
		case "allMessage":
			break
		default:
			for _, client := range users {
				if err = client.Websocket.WriteMessage(messageType, message); err != nil {
					log.Println("Cloud not send Message to ", client.ClientIP, err.Error())
				}
			}
		}
	}
}

func init() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/ws", WebsocketHandler)
}

func main() {
	log.Fatal(http.ListenAndServe("localhost:8002", nil))
}

