import React from 'react'

import '~/css/sidebar.scss'

class Sidebar extends React.Component {
    render() {
        const element = (
            <div className="sidebar">
                <div className="user-avatar">
                    <img src="/public/images/avatar-example.jpg" />
                </div>
                <i className="iconfont iconliaotian active"></i>
                <i className="iconfont iconlianxiren"></i>
                <i className="iconfont iconicon-test"></i>
                <i className="iconfont iconshezhi1"></i>
            </div>
        )
        return element
    }
}

module.exports = Sidebar