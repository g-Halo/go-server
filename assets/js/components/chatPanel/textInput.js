import React from 'react'

class TextInput extends React.Component {
  render() {
    return (
        <div className="chat-text-input">
            <textarea placeholder="enter 键发送"></textarea>
        </div>
    )
  }
}

module.exports = TextInput