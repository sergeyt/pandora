import React from "react";
import {ListSubheader} from "@material-ui/core";
import List from "@material-ui/core/List";
import PropTypes from "prop-types";

export function SideBarCategory(props) {
    const {title, children, ...rest} = props;
    const subheader = title ? (<ListSubheader inset>{title}</ListSubheader>) : null;

    return (
        <List {...rest}>
            <div>
                {subheader}
                {children}
            </div>
        </List>
    );
}

SideBarCategory.propTypes = {
    title: PropTypes.string,
    children: PropTypes.oneOfType([
        PropTypes.arrayOf(PropTypes.node),
        PropTypes.node
    ]).isRequired
};

export default SideBarCategory;
