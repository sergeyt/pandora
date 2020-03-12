import React, {useState} from "react";
import {makeStyles} from "@material-ui/styles";
import {AppBar, Container, Drawer, IconButton, Toolbar, Typography} from "@material-ui/core";
import MenuIcon from "@material-ui/icons/Menu";
import PropTypes from "prop-types";
import SideBar from "./SideBar";

const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
    root: {
        display: "flex",
    },
    appBar: {
        marginLeft: drawerWidth,
        zIndex: theme.zIndex.drawer + 1,
        width: "100%",
        [theme.breakpoints.down("sm")]: {
            marginLeft: 0,
        },
    },
    menuButton: {
        marginRight: theme.spacing(2),
        [theme.breakpoints.up("md")]: {
            display: "none",
        },
    },
    title: {
        flexGrow: 1,
    },
    permanentDrawerPaper: {
        // without this property content will
        // be hidden behind the drawer.
        position: "relative",
        width: drawerWidth,
        [theme.breakpoints.down("sm")]: {
            display: "none",
        }
    },
    temporaryDrawerPaper: {
        width: drawerWidth,
    },
    content: {
        flexGrow: 1,
        height: "100vh",
        overflow: "auto",
    },
    contentSpacer: theme.mixins.toolbar,
    contentContainer: {
        paddingTop: theme.spacing(2),
        paddingBottom: theme.spacing(2),
        height: '100%',
    },
}));

function Dashboard(props) {
    const {sidebar, content} = props;

    const classes = useStyles();
    const [menuOpen, setMenuOpen] = useState(false);
    const openMenu = () => setMenuOpen(true);
    const closeMenu = () => setMenuOpen(false);

    return (
        <div className={classes.root}>
            <AppBar className={classes.appBar}>
                <Toolbar>
                    <IconButton
                        edge="start"
                        color="inherit"
                        aria-label="open drawer"
                        onClick={openMenu}
                        className={classes.menuButton}
                    >
                        <MenuIcon/>
                    </IconButton>
                    <Typography component="h1" variant="h6" color="inherit" noWrap className={classes.title}>
                        Pandora
                    </Typography>
                </Toolbar>
            </AppBar>
            <Drawer variant="permanent" classes={{paper: classes.permanentDrawerPaper}}>
                <SideBar>{sidebar}</SideBar>
            </Drawer>
            <Drawer open={menuOpen} classes={{paper: classes.temporaryDrawerPaper}} onClose={closeMenu}>
                <SideBar variant="button" onClose={closeMenu}>{sidebar}</SideBar>
            </Drawer>
            <main className={classes.content}>
                <div className={classes.contentSpacer}/>
                <Container className={classes.contentContainer}>
                    {content}
                </Container>
            </main>
        </div>
    );
}

Dashboard.propTypes = {
    content: PropTypes.oneOfType([
        PropTypes.arrayOf(PropTypes.node),
        PropTypes.node
    ]).isRequired,
    sidebar: PropTypes.oneOfType([
        PropTypes.arrayOf(PropTypes.node),
        PropTypes.node
    ]).isRequired
};

export default Dashboard;
