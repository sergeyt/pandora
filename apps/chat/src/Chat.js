import React from 'react'
import { createSelector } from 'reselect'

import { connect } from './connect'
import Message from './Message'

class Input extends React.Component {
  state = {
    value: '',
  }

  onKeyUp = e => {
    if (e.which === 13) {
    }
  }

  render() {
    return (
      <input
        type="text"
        value={this.state.value}
        onChange={e =>
          this.setState({
            value: e.target.value,
          })
        }
        onKeyUp={this.onKeyUp}
      />
    )
  }
}

class Chat extends React.Component {
  render() {
    const content = this.props.messages.map((m, i) => (
      <Message key={i} msg={m} />
    ))
    return (
      <div className="chat">
        <div className="messages">{content}</div>
        <div className="input-container">
          <Input />
        </div>
      </div>
    )
  }
}

const messages = createSelector(
  state => state.chat.messages,
  messages => ({ messages })
)

export default connect({ messages })(Chat)
