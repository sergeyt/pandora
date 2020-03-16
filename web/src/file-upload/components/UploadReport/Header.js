import React from "react";
import PropTypes from "prop-types";
import CollapseIcon from "@material-ui/icons/ExpandMore";
import ExpandIcon from "@material-ui/icons/ExpandLess";
import CloseIcon from "@material-ui/icons/Close";
import IconButton from "@material-ui/core/IconButton";
import SnackbarContent from "@material-ui/core/SnackbarContent";
import Button from "@material-ui/core/Button";


function Header(props) {
    const {
        count,
        done,
        expanded,
        onCancel,
        onExpand,
        onClose,
    } = props;

    const ExpandButtonIcon = (expanded ? CollapseIcon : ExpandIcon);

    return (
        <SnackbarContent
            message={`${count} uploads`}
            action={
                <React.Fragment>
                    {
                        !done && (
                            <Button color="secondary" size="small" onClick={onCancel}>
                                Cancel
                            </Button>
                        )
                    }
                    <IconButton size="small" aria-label="delete" onClick={() => onExpand(!expanded)}>
                        <ExpandButtonIcon
                            fontSize="small"
                            color="secondary"
                        />
                    </IconButton>
                    <IconButton size="small" aria-label="delete" onClick={onClose}>
                        <CloseIcon fontSize="small" color="secondary"/>
                    </IconButton>
                </React.Fragment>
            }
        />
    );
}

Header.propTypes = {
    count: PropTypes.number.isRequired,
    done: PropTypes.bool.isRequired,
    expanded: PropTypes.bool.isRequired,
    onCancel: PropTypes.func.isRequired,
    onExpand: PropTypes.func.isRequired,
    onClose: PropTypes.func.isRequired,
};

export default Header;
