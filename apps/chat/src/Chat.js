import React from 'react'

import './Chat.css'
import { connect } from './connect'
import Message from './Message'
import ChatInput from './ChatInput'
import { currentUser, messages } from './selectors'
import { send } from './saga'

class Chat extends React.Component {
  processMessage = content => {
    // TODO process commands
    const message = {
      from: this.props.currentUser.uid,
      // TODO specify channel id
      to: 'general',
      content,
    }
    this.props.send(message)
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

export default connect(
  { currentUser, messages },
  { send }
)(Chat)
