import React from 'react';
import ReactDOM from 'react-dom';

import '~/css/main.scss'


import Sidebar from "./components/sidebar"
import ContactPanel from './components/contactPanel'
import ChatPanel from './components/chatPanel'
const element = (
    <div className="chat-container">
        <Sidebar></Sidebar>
        <ContactPanel></ContactPanel>
        <ChatPanel></ChatPanel>
    </div>
);

ReactDOM.render(
    element,
    document.getElementById('chat')
);