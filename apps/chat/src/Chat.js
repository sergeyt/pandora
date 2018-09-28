import React from 'react';

import Message from './Message';

export default class Chat extends React.Component {
    render() {
        const content = this.props.messages.map((m, i) => <Message key={i} msg={m}/>);
        return (
            <div class="chat">
                <div class="messages">
                    {content}
                </div>
                <div class="input-container">
                </div>
            </div>
        );
    }
}
