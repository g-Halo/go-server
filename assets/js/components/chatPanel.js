import React from 'react'
import Header from './chat-panel/header'
import MessagePanel from './chat-panel/messagePanel'
import TextHeader from './chat-panel/textHeader'
import TextInput from './chat-panel/textInput'
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