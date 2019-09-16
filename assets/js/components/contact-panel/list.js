import React from 'react'
import { Get } from 'react-axios'
import '~/css/contact-panel.scss'

class List extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            users: [
                { id: 1, name: 'Nancy', desc: '在吗？', time: '3 分钟前', unread: 2 },
                { id: 2, name: 'Cat', desc: '上次吃饭没给你钱呢', time: '4 分钟前', unread: 1 },
                { id: 3, name: 'Macp', desc: '不说了，洗澡了', time: '4 分钟前', unread: 0 },
                { id: 4, name: 'Lucy', desc: '来啊，打王者啊', time: '10 分钟前', unread: 0 },
                { id: 5, name: 'Cury', desc: '今天玩得开心吗', time: '1 天前', unread: 20 },
                { id: 6, name: 'Yanx', desc: '[动画表情]', time: '1 天前', unread: 10 },
                { id: 7, name: 'Unix', desc: '对方已领取红包', time: '2 天前', unread: 1 },
                { id: 8, name: 'Pos', desc: '视频聊天[已结束]', time: '2 天前', unread: 0 },
                { id: 9, name: 'caw3', desc: '语音视频[已结束]', time: '2 天前', unread: 0 },
                { id: 10, name: '萨科', desc: '你好', time: '8 天前', unread: 0 }
            ],
            activeUserId: 2
        }
    }
    
    onSwitchUser(user) {
        this.setState({
            activeUserId: user.username
        })
    }

    render() {
        const userInfo = (
            <Get url="/v1/contacts">
                {(error, response) => {
                   if (response !== null) {
                       const users = response.data
                       return users.map((user) =>
                            <div
                                key={user.username}
                                className={`contact-panel__user ${this.state.activeUserId === user.id ? 'active' : ''}`}
                                onClick={this.onSwitchUser.bind(this, user)}
                            >
                                <div className="user-avatar">
                                    <img src="/public/images/avatar-example.jpg" />
                                </div>
                                <div className="contact-panel-user__info">
                                    <div className="user__header">
                                        <span className="user__name">{user.nickname}</span>
                                        <span className="text-gray fs12">{user.last_chat_at}</span>
                                    </div>
                                    <div className="user__desc">
                                        <span className="text-gray message-desc">{user.desc}</span>
                                        { user.unread > 0 ? (<span className="unread-circle">{user.unread}</span>) : '' }
                                    </div>
                                </div>
                            </div>
                        )
                   } else {
                       return (<div>Get Error</div>)
                   }
                }}
            </Get>
        )

        return (
            <div className="contact-panel__list">
                { userInfo }
            </div>
        )
    }
}

module.exports = List