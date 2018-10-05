import React from 'react'
import { createSelector } from 'reselect'

import './Chat.css'
import { connect } from './connect'
import Message from './Message'
import ChatInput from './ChatInput'
import { pushMessage } from './actions'

class Chat extends React.Component {
  processMessage = content => {
    const message = {
      content,
    }
    this.props.pushMessage(message)
  }

  render() {
    const content = this.props.messages.map((m, i) => (
      <Message key={i} message={m} />
    ))
    return (
      <div className="chat-container">
        <ChatInput onMessage={this.processMessage} />
        <div className="chat">
          <div style={{ height: '10px' }} />
          {content}
          <div style={{ height: '85px' }} />
        </div>
      </div>
    )
  }
}

const messages = createSelector(
  state => state.chat.messages,
  messages => ({ messages })
)

export default connect(
  { messages },
  { pushMessage }
)(Chat)
