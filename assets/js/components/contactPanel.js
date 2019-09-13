import React from 'react'
import Header from './contact-panel/header'
import Search from './contact-panel/search'
import List from './contact-panel/list'
class ContactPanel extends React.Component {
    render() {
        const e = (
            <div className="contact-panel">
                <Header></Header>
                <Search></Search>
                <List></List>
            </div>
        )
        return e
    }
}

module.exports = ContactPanel