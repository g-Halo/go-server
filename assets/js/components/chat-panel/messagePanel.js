import React from 'react'

class MessagePanel extends React.Component {
  render() {
    return (
        <div className="chat-message-panel">
            <div className="chat-message__time fs12">
              <span>13:03</span>
            </div>

            <div className="target-user">
              <img className="chat avatar small" src="/public/images/avatar-example.jpg" />
              <div className="message-box">
                Hi Joy, How are you
              </div>
            </div>

            <div className="owner">
              <div className="message-box">
                I am fine.Thank you.
              </div>
              <img className="chat avatar small" src="/public/images/avatar-example.jpg" />
            </div>

            <div className="owner">
              <div className="message-box">
                Long Time no see.
              </div>
              <img className="chat avatar small" src="/public/images/avatar-example.jpg" />
            </div>
            
            <div className="target-user">
              <img className="chat avatar small" src="/public/images/avatar-example.jpg" />
              <div className="message-box">
                呵呵
              </div>
            </div>
        </div>
    )
  }
}

module.exports = MessagePanel