import React from 'react'

class Header extends React.Component {
  render() {
    return (
        <div className="chat-header">
            <div className="chat-header__user-info">
                <img className="chat avatar normal" src="/public/images/avatar-example.jpg" />
                <div className="username">
                    <div>Nancy</div>
                    <div className="text-gray fs12">在线</div>
                </div>
            </div>
            <div className="header-operate">
                <i className="iconfont iconicon_tianjiahaoyou"></i>
                <i className="iconfont iconjilu"></i>
            </div>
        </div>
    )
  }
}

module.exports = Header