import React from 'react'

class Header extends React.Component {
  render() {
    return (
        <div className="chat-header">
            <img className="chat avatar normal" src="/public/images/avatar-example.jpg" />
            <div className="username">
                <div>Nancy</div>
                <div className="text-gray fs12">在线</div>
            </div>
        </div>
    )
  }
}

module.exports = Header