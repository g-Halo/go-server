import React from 'react'
import Header from './contactPanel/header'
import Search from './contactPanel/search'
import List from './contactPanel/list'
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