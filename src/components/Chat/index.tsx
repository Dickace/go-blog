import React, {ChangeEvent, useEffect, useState} from 'react'
import {useHistory} from "react-router-dom";
import {SCREENS} from "../../routes/endpoints";
import URLS from '../../ApiUrl.json'
import SendIcon from '@mui/icons-material/Send';
import {Button, List, ListItem, ListItemButton, ListItemText, TextField} from "@mui/material";
import './style.css'

type User = {
    username: string
    uid: string
}


type Message = {
    from: string
    to: string
    text: string
}



const Chat : React.FC = () => {
    const history = useHistory()

    const username:string |null= localStorage.getItem("username")
    const [userList,setUserList] = useState<Array<User>>([])
    const [messages,setMessages] = useState<Array<Message>>([])
    const [myProfile,setMyProfile] = useState<User>({}as User)
    const [currentChat, setCurrentChat] = useState<Array<Message>>([])
    const [selectedListItem,setSelectedListItem] = useState<number>()
    const [currentUser, setCurrentUser] = useState<string>("")
    const [message,setMessage] = useState<string>("")
    if( username  === ""  || username === null){
        history.push(SCREENS.login_screen)
    }
    const [webSocket,setWebSocket] = useState<WebSocket|undefined>(undefined)
    const handleAllChatClick = (event: React.MouseEvent<HTMLDivElement, MouseEvent>,index: number) => {
        handleListItemClick(event, index)
        SetAllMessages()
        setCurrentUser("")
    }
    const handleUserClick = (event: React.MouseEvent<HTMLDivElement, MouseEvent>,index: number) =>{

        const uidOtherUser =  event.currentTarget.getAttribute("data-uid")
        if (uidOtherUser!== null){
            SetMessages(uidOtherUser)
            setCurrentUser(uidOtherUser)
        }
        handleListItemClick(event, index)
    }
    const handleListItemClick = (
        event: React.MouseEvent<HTMLDivElement, MouseEvent>,
        index: number,
    ) => {
        setSelectedListItem(index);
    };
    const SetMessages = (user: string) =>{
        const msgs : Array<Message> = []
        messages.map((value => {
            if ((value.to === myProfile.uid && value.from === user)||(value.from === myProfile.uid && value.to=== user)){
                msgs.push(value)
            }
        }))
        setCurrentChat(msgs)
    }
    const handleChangeMessage = (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>)  => {
        setMessage(event.currentTarget.value)
    }
    const handleSendMessage = () =>{

        if (webSocket!==undefined){
            let msg : Message ={
                from: myProfile.uid,
                text: message.toString(),
                to: currentUser
            }
            console.log(msg)
            if(currentUser === ""){
                webSocket.send( JSON.stringify({"type":"allMessage","text":`"to":"${msg.to}","from":"${msg.from}", "text":"${msg.text}"`}))
            } else {
                webSocket.send( JSON.stringify({"type":"message","text":`"to":"${msg.to}","from":"${msg.from}", "text":"${msg.text}"`}))
            }
            setMessages([...messages, msg])

        }
    }
    const SetAllMessages = () =>{
        const msgs : Array<Message> = []
        messages.map((value => {
            if (value.to === ""){
                msgs.push(value)
            }
        }))
        setCurrentChat(msgs)
    }
    useEffect(()=>{
        localStorage.setItem("chats",JSON.stringify(messages))
        if (currentUser !== ""){
            SetMessages(currentUser)
        } else {
            SetAllMessages()
        }

    },[messages])
    useEffect(
        ()=>{
            let websocket =new WebSocket(
                `${URLS.WS_URL}?username=${username}`
            )
            setWebSocket(websocket)
            let encodeMsgs = localStorage.getItem("chats")
            console.log(encodeMsgs)
            if (encodeMsgs!== null){
                setMessages(JSON.parse(encodeMsgs))
            }
            websocket.onmessage = async function(msg){
                try {
                    const msgData = JSON.parse(msg.data)

                    switch (msgData?.type) {
                        case "newClient":
                            console.log(userList)
                            let newUser = {username: msgData.users[0].username, uid: msgData.users[0].uid} as User
                            setUserList((prevState)=>[...prevState, newUser])
                            break
                        case "userList":
                            let newUsers : Array<User>=[]
                            let users = msgData.users
                            if (users!==null){
                                users.map((value: { username: string; uid: string; })=>{
                                    newUsers.push({username: value.username,
                                        uid: value.uid
                                    } as User)
                                })
                            }
                            setUserList(
                                newUsers
                            )
                            break
                        case "message":
                            console.log(messages)
                            setMessages((prevState)=>[...prevState, {
                                from: msgData.message.from,
                                text: msgData.message.text,
                                to: msgData.message.to
                            } as Message])
                            break
                        case "allMessage":
                            setMessages((prevState)=>[...prevState, {
                                from: msgData.message.from,
                                text: msgData.message.text,
                                to: ""
                            } as Message])
                            break
                        case "myData":
                            setMyProfile({username: msgData.users[0].username,
                                            uid: msgData.users[0].uid
                            } as User)
                            break
                    }
                }
                catch (e){
                    console.log(e)
                }
            }
            websocket.onopen = function () {
                websocket.send(JSON.stringify({ type: 'userList' }))
            }
        },[]


    )

    return(
        <div className="chat">
            <div className="userList">
                <List>
                        <ListItemButton selected={selectedListItem === 0}  onClick={(event)=>{handleAllChatClick(event,0)}}>
                            <ListItemText primary="Общий чат"/>
                        </ListItemButton>

                    {userList.map((value,index)=>{
                        return(
                            <ListItemButton selected={selectedListItem === index+1} onClick={(event)=>{handleUserClick(event,index+1)}} data-uid={value.uid}>
                               <ListItemText primary={value.username}/>
                            </ListItemButton>
                        )

                    })}
                </List>
            </div>

            <div className="chatArea">
                <div className="messageBox">
                    {currentChat.map((value)=>{
                        console.log(value)
                        let fromMe = false
                        if (value.from === myProfile.uid){
                            fromMe = true
                        }
                        return(<div className={`messageBlock ${fromMe? "messageBlock_fromMe" : ""}`} ><p>{value.text} </p></div>)
                    })}
                </div>
                <div className="chatArea-input">

                    <TextField value={message} onChange={handleChangeMessage} />
                    <Button onClick={handleSendMessage}>
                        <SendIcon/>
                    </Button>
                </div>
            </div>

        </div>
    )
}

export default Chat