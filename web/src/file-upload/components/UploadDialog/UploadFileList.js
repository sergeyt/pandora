import React from "react";
import PropTypes from "prop-types";
import {makeStyles} from "@material-ui/styles";
import Chip from "@material-ui/core/Chip";

const useStyles = makeStyles(theme => ({
    container: {
        display: "flex",
        flexDirection: "row",
        alignItems: "flex-start",
        flexWrap: "wrap",
    },
    file: {
        margin: theme.spacing(0.5),
    }
}));

function UploadFileList(props) {
    const {
        files,
        onDelete,
    } = props;

    const classes = useStyles();
    return (
        <div className={classes.container}>
            {
                files.map((file, i) => (
                    <Chip
                        key={file.path}
                        label={file.name}
                        onDelete={() => onDelete(file, i)}
                        className={classes.file}
                    />
                ))
            }
        </div>
    );
}

UploadFileList.propTypes = {
    files: PropTypes.arrayOf(PropTypes.object).isRequired,
    onDelete: PropTypes.func.isRequired,
};

export default UploadFileList;