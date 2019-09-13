import React from 'react'

class Search extends React.Component {
    render() {
        return (
            <div className="search-input">
                <input type="text" placeholder="搜索用户" />
            </div>
        )
    }
}

module.exports = Search