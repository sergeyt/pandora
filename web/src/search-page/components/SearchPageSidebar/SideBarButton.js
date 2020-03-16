import React from "react";
import PropTypes from "prop-types";
import ListItem from "@material-ui/core/ListItem";
import {ListItemIcon, ListItemText} from "@material-ui/core";


function SideBarButton(props) {
    const {
        icon: Icon,
        text,
        onClick
    } = props;

    return (
        <ListItem button onClick={onClick}>
            <ListItemIcon>
                <Icon/>
            </ListItemIcon>
            <ListItemText primary={text}/>
        </ListItem>
    );
}

SideBarButton.propTypes = {
    icon: PropTypes.elementType.isRequired,
    text: PropTypes.string.isRequired,
    onClick: PropTypes.func,
};

export default SideBarButton;
