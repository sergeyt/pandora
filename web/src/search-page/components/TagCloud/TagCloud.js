import React from "react";
import PropTypes from "prop-types";
import {makeStyles} from "@material-ui/core/styles";
import Chip from "@material-ui/core/Chip";

const useStyles = makeStyles(theme => ({
    container: {
        display: "flex",
        flexWrap: "wrap",
        paddingRight: theme.spacing(2),
        paddingLeft: theme.spacing(2),
    },
    tag: {
        margin: theme.spacing(0.5),
    }
}));

function TagCloud(props) {
    const {
        tags,
        onClick
    } = props;

    const classes = useStyles();

    return (
        <div className={classes.container}>
            {
                tags.map(tag => (
                    <Chip
                        key={tag}
                        label={tag}
                        onClick={() => onClick(tag)}
                        className={classes.tag}
                    />
                ))
            }
        </div>
    );
}

TagCloud.propTypes = {
    tags: PropTypes.arrayOf(PropTypes.string).isRequired,
    onClick: PropTypes.func.isRequired,
};

export default TagCloud;