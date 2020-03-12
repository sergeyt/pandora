import React from "react";
import PropTypes from "prop-types";
import {makeStyles} from "@material-ui/core/styles";
import Card from "@material-ui/core/Card";
import CardHeader from "@material-ui/core/CardHeader";
import CardMedia from "@material-ui/core/CardMedia";
import CardContent from "@material-ui/core/CardContent";
import CardActions from "@material-ui/core/CardActions";
import Avatar from "@material-ui/core/Avatar";
import IconButton from "@material-ui/core/IconButton";
import Typography from "@material-ui/core/Typography";
import {red} from "@material-ui/core/colors";
import FavoriteIcon from "@material-ui/icons/Favorite";
import MoreVertIcon from "@material-ui/icons/MoreVert";
import Chip from "@material-ui/core/Chip";

const useStyles = makeStyles(theme => ({
    root: {

    },
    media: {
        height: 0,
        paddingTop: "56.25%", // 16:9
    },
    expand: {
        transform: "rotate(0deg)",
        marginLeft: "auto",
        transition: theme.transitions.create("transform", {
            duration: theme.transitions.duration.shortest,
        }),
    },
    expandOpen: {
        transform: "rotate(180deg)",
    },
    avatar: {
        backgroundColor: red[500],
    },
    tags: {
        display: "flex",
        justifyContent: "flex-start",
        flexWrap: "wrap",
        "& > *": {
            margin: theme.spacing(0.5),
        },
    },
}));

function DocumentPreview(props) {
    const {document} = props;
    const classes = useStyles();

    return (
        <Card className={classes.root}>
            <CardHeader
                avatar={
                    <Avatar aria-label="recipe" className={classes.avatar}>
                        PDF
                    </Avatar>
                }
                action={
                    <IconButton aria-label="settings">
                        <MoreVertIcon/>
                    </IconButton>
                }
                title={document.title}
                subheader={`${document.date} - ${document.size}`}
            />
            <CardMedia
                className={classes.media}
                image={document.image}
                title={document.title}
            />
            <CardContent>

                <div className={classes.tags}>
                    {document.tags.map(tag => (
                        <Chip
                            variant="outlined"
                            label={tag}
                            onClick={() => null}
                            key={tag}
                        />
                    ))}
                </div>
                <Typography variant="body2" color="textSecondary" component="p">
                    {document.previewText}
                </Typography>
            </CardContent>
            <CardActions disableSpacing>
                <IconButton aria-label="add to favorites">
                    <FavoriteIcon/>
                </IconButton>
            </CardActions>
        </Card>
    );
}

DocumentPreview.propTypes = {
    document: PropTypes.shape({
        title: PropTypes.string.isRequired,
        date: PropTypes.string.isRequired,
        size: PropTypes.string.isRequired,
        image: PropTypes.string.isRequired,
        tags: PropTypes.arrayOf(PropTypes.string).isRequired,
        link: PropTypes.string.isRequired,
        previewText: PropTypes.string,
    })
};


export default DocumentPreview;