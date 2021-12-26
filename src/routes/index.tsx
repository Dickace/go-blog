import React from 'react'
import { Switch, Route, Redirect, withRouter } from 'react-router-dom'
import {SCREENS} from "./endpoints";
import Auth from "../components/Auth";
import Chat from "../components/Chat";

const Routes: React.FC = () => {
    return (
        <Switch>
            <Route path={SCREENS.chat_screen} >
                <Chat/>
            </Route>
            <Route path={SCREENS.login_screen}>
                <Auth/>
            </Route>
            <Redirect to={SCREENS.login_screen} />
        </Switch>
    )
}

export default withRouter(Routes)