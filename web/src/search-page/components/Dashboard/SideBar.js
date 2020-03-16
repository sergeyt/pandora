import React from "react";
import PropTypes from "prop-types";
import {makeStyles} from "@material-ui/styles";
import {IconButton} from "@material-ui/core";
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";


const useStyles = makeStyles(theme => ({
    spacer: theme.mixins.toolbar,
    closeButton: {
        display: "flex",
        alignItems: "center",
        justifyContent: "flex-end",
        padding: "0 8px",
        ...theme.mixins.toolbar,
    },
}));

const VARIANT_NONE = "none";
const VARIANT_SPACER = "spacer";
const VARIANT_BUTTON = "button";

function SideBar(props) {
    const classes = useStyles();
    const {variant = VARIANT_SPACER, onClose, children} = props;

    let firstItem = null;
    if (variant === VARIANT_SPACER) {
        firstItem = <div className={classes.spacer}/>;
    } else if (variant === VARIANT_BUTTON) {
        firstItem = (
            <div className={classes.closeButton}>
                <IconButton onClick={onClose}>
                    <ChevronLeftIcon/>
                </IconButton>
            </div>
        );
    }

    return (
        <React.Fragment>
            {firstItem}
            {children}
        </React.Fragment>
    );
}

SideBar.propTypes = {
    onClose: PropTypes.func,
    variant: PropTypes.oneOf([VARIANT_NONE, VARIANT_SPACER, VARIANT_BUTTON]),
    children: PropTypes.oneOfType([
        PropTypes.arrayOf(PropTypes.node),
        PropTypes.node
    ]).isRequired
};

export default SideBar;
