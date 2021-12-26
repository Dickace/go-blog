import React, {ChangeEvent, useState} from "react";
import {Box, TextField, Button} from "@mui/material";
import {AccountCircle} from "@mui/icons-material";
import './style.css'
import {useHistory} from "react-router-dom";
import {SCREENS} from "../../routes/endpoints";

const Auth: React.FC = () => {
    const history = useHistory();
    const [username, setUsername] = useState<string>("")
    const handleLoginClick = () => {
        localStorage.setItem("username", username)
        history.push(`${SCREENS.chat_screen}`)
    }
    const handleChangeUsername = (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        setUsername(event.currentTarget.value)
    }

    return (
        <div className="loginPage">
            <Box>
                <Box sx={{display: 'flex', alignItems: 'flex-end'}}>
                    <AccountCircle sx={{color: 'action.active', mr: 1, my: 0.5}}/>
                    <TextField value={username} onChange={handleChangeUsername} id="input-name" label="Enter your name"
                               variant="standard"/>
                </Box>

                <Button onClick={handleLoginClick} variant="contained">Login to chat</Button>
            </Box>

        </div>
    )


}

export default Auth