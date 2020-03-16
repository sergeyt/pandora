import React, {useState} from "react";
import PropTypes from "prop-types";
import ListItem from "@material-ui/core/ListItem";
import FileIcon from "@material-ui/icons/InsertDriveFile";
import SuccessIcon from "@material-ui/icons/CheckCircle";
import CancelIcon from "@material-ui/icons/Cancel";
import ErrorIcon from "@material-ui/icons/Error";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import ListItemSecondaryAction from "@material-ui/core/ListItemSecondaryAction";
import IconButton from "@material-ui/core/IconButton";
import CircularProgress from "@material-ui/core/CircularProgress";
import {makeStyles} from "@material-ui/core/styles";


const useStyles = makeStyles(theme => ({
    success: {
        color: theme.palette.success.main,
    },
    error: {
        color: theme.palette.error.main,
    }
}));

const maxNameLength = 20;

function displayName(path) {
    if (path.length < maxNameLength) {
        return path;
    }
    return path.slice(0, maxNameLength - 3) + "...";
}

function ItemButton(props) {
    const {
        hover,
        status,
        onCancel,
    } = props;

    const classes = useStyles();

    if (status === "success") {
        return (<SuccessIcon className={classes.success}/>);
    } else if (status === "failure") {
        return (<ErrorIcon className={classes.error}/>);
    } else if (status === "active" && !hover) {
        return (<CircularProgress size={20}/>);
    } else {
        return (
            <IconButton onClick={onCancel} edge="end">
                <CancelIcon/>
            </IconButton>
        );
    }
}

ItemButton.propTypes = {
    hover: PropTypes.bool.isRequired,
    status: PropTypes.oneOf(["pending", "active", "success", "failure"]).isRequired,
    onCancel: PropTypes.func.isRequired,
};


function FileListItem(props) {
    const {
        file,
        onCancel
    } = props;

    const [hover, setHover] = useState(false);


    return (
        <ListItem onMouseEnter={() => setHover(true)} onMouseLeave={() => setHover(false)}>
            <ListItemIcon>
                <FileIcon color="action"/>
            </ListItemIcon>
            <ListItemText
                primary={displayName(file.path)}
                secondary={file.status}
            />
            <ListItemSecondaryAction onMouseEnter={() => setHover(true)} onMouseLeave={() => setHover(false)}>
                <ItemButton status={file.status} hover={hover} onCancel={() => onCancel(file)}/>
            </ListItemSecondaryAction>
        </ListItem>
    );
}

FileListItem.propTypes = {
    file: PropTypes.shape({
        path: PropTypes.string.isRequired,
        status: PropTypes.oneOf(["pending", "active", "success", "failure"]).isRequired,
    }).isRequired,
    onCancel: PropTypes.func.isRequired,
};

export default FileListItem;
