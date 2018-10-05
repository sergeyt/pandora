import React from 'react'

class ChatInput extends React.Component {
  state = {
    value: '',
  }

  onKeyUp = e => {
    const key = e.which || e.keyCode
    var value = e.target.value.replace(/^(\s)|(\s+)$/g, '') // removes whitespace before and after
    if (key === 13 && value) {
      this.setState({
        value: '',
      })
      this.props.onMessage(value)
    }
  }

  render() {
    return (
      <input
        type="text"
        className="chat-input dropIn"
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

export default ChatInput
