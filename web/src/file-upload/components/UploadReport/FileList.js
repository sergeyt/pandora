import React from "react";
import PropTypes from "prop-types";
import List from "@material-ui/core/List";
import {makeStyles} from "@material-ui/core/styles";
import FileListItem from "./FileListItem";
import Divider from "@material-ui/core/Divider";

const useStyles = makeStyles(() => ({
    list: {
        maxHeight: "300px",
        width: "100%",
        overflow: "auto"
    },
}));

function FileList(props) {
    const {
        files,
        onCancel,
    } = props;

    const classes = useStyles();

    return (
        <List dense className={classes.list}>
            {files.map((file, i) => (
                <React.Fragment key={file.path}>
                    <FileListItem
                        file={file}
                        onCancel={onCancel}
                    />
                    {(i !== files.length - 1 ? <Divider/> : null)}
                </React.Fragment>
            ))}
        </List>
    );
}

FileList.propTypes = {
    files: PropTypes.arrayOf(
        PropTypes.shape({
            path: PropTypes.string.isRequired,
            status: PropTypes.oneOf(["pending", "active", "success", "failure"]).isRequired,
        })
    ).isRequired,
    onCancel: PropTypes.func.isRequired,
};

export default FileList;
