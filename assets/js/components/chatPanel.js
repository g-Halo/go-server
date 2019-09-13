import React from 'react'
import Header from './chatPanel/header'
import MessagePanel from './chatPanel/messagePanel'
import TextHeader from './chatPanel/textHeader'
import TextInput from './chatPanel/textInput'
import '~/css/chat-panel.scss'

class ChatPanel extends React.Component {
  render() {
    return (
      <div className="chat-panel">
        <Header></Header>
        <MessagePanel></MessagePanel>
        <TextHeader></TextHeader>
        <TextInput></TextInput>
      </div>
    )
  }
}

module.exports = ChatPanel