import React from 'react'

export default class Message extends React.Component {
  render() {
    return (
      <div className="chat-msg general show">{this.props.message.content}</div>
    )
  }
}
