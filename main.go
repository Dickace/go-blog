package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
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
	ClientIP  string
	UserData  User
}

type User struct {
	Username string
	Messages []Message
	UID      string
}

type MessageStruct struct {
	To   string `json:"to"`
	From string `json:"from"`
	Text string `json:"text"`
}

type MessageSend struct {
	Type    string        `json:"type"`
	Message MessageStruct `json:"message"`
}

type Message struct {
	Types string `json:"type"`
	Text  string `json:"text"`
}
type MsgToClient struct {
	Type  string         `json:"type"`
	Users []UserResponse `json:"users"`
}

type UserResponse struct {
	Username string `json:"username"`
	Uid      string `json:"uid"`
}

func newConnectUser(ws *websocket.Conn, clientIP string, r *http.Request) *ConnectUser {
	u, err := uuid.NewV4()
	if err != nil {
		fmt.Println("Can't generate uuid")
	}
	return &ConnectUser{
		Websocket: ws,
		ClientIP:  clientIP,
		UserData: User{
			Username: r.URL.Query().Get("username"),
			Messages: []Message{},
			UID:      u.String(),
		},
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var users = []ConnectUser{}

func RemoveIndex(s []ConnectUser, index int) []ConnectUser {
	return append(s[:index], s[index+1:]...)
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, _ := upgrader.Upgrade(w, r, nil)

	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("Websocket could not be closed", err.Error())
		}
	}()

	log.Println("Client connected:", ws.RemoteAddr().String())
	var socketClient *ConnectUser = newConnectUser(ws, ws.RemoteAddr().String(), r)
	users = append(users, *socketClient)
	log.Println("Number client connected ...", len(users))

	var msgToClient = MsgToClient{
		Type: "newClient",
		Users: []UserResponse{
			{Username: socketClient.UserData.Username,
				Uid: socketClient.UserData.UID},
		},
	}
	respMsg, err := json.Marshal(msgToClient)
	if err != nil {
		fmt.Print("Cannot encode Json")
	}
	for _, client := range users {
		if client.UserData.UID != socketClient.UserData.UID {
			if err = client.Websocket.WriteMessage(1, respMsg); err != nil {
				log.Println("Cloud not send Message to ", client.ClientIP, err.Error())
			}
		}
	}

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Ws disconnect waiting", err.Error())
			for i, v := range users {
				if v.ClientIP == socketClient.ClientIP {
					users = RemoveIndex(users, i)
				}
			}
			log.Println("Number of client still connected ...", len(users))
			return
		}

		var msg *Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Print("Decode error")
		}
		switch msg.Types {
		case "userList":
			var userLists []UserResponse
			for _, v := range users {
				userLists = append(userLists, UserResponse{
					Username: v.UserData.Username,
					Uid:      v.UserData.UID,
				})
			}
			var msgToUser = MsgToClient{
				Type:  "userList",
				Users: userLists,
			}
			respMsg, err := json.Marshal(msgToUser)
			if err != nil {
				fmt.Print("Cannot encode Json")
			}
			if err = socketClient.Websocket.WriteMessage(messageType, respMsg); err != nil {
				log.Println("Cannot send UserList")
			}
			break
		case "message":

			var msgDecode *MessageStruct
			err = json.Unmarshal([]byte("{"+msg.Text+"}"), &msgDecode)
			if err != nil {
				fmt.Println("Decode json error")
			}
			for _, v := range users {
				if v.UserData.UID == msgDecode.To {
					var msgToClient = MessageSend{
						Type: "message",
						Message: MessageStruct{
							From: socketClient.UserData.UID,
							Text: msgDecode.Text,
							To:   msgDecode.To,
						},
					}
					msgs, err := json.Marshal(msgToClient)
					if err != nil {
						fmt.Println("Error encode Json")
					}
					if err = v.Websocket.WriteMessage(messageType, msgs); err != nil {
						log.Println("Cannot send UserList")
					}
				}
			}
			break
		case "allMessage":
			var msgDecode *MessageStruct
			fmt.Print(msg.Text)
			err = json.Unmarshal([]byte("{"+msg.Text+"}"), &msgDecode)
			if err != nil {
				fmt.Println("Decode json error")
			}

			for _, client := range users {
				if client.UserData.UID != socketClient.UserData.UID {
					var msgToClient = MessageSend{
						Type: "message",
						Message: MessageStruct{
							From: socketClient.UserData.UID,
							Text: msgDecode.Text,
							To:   client.UserData.UID,
						},
					}
					msgs, err := json.Marshal(msgToClient)
					if err != nil {
						fmt.Println("Error encode Json")
					}
					if err = client.Websocket.WriteMessage(1, msgs); err != nil {
						log.Println("Cloud not send Message to ", client.ClientIP, err.Error())
					}
				}
			}
			break
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
