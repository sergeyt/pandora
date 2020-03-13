import React from "react";
import PropTypes from "prop-types";
import SideBarCategory from "./SideBarCategory";
import SideBarButton from "./SideBarButton";
import SearchIcon from "@material-ui/icons/Search";
import UploadIcon from "@material-ui/icons/CloudUpload";
import TagCloud from "../TagCloud";

function SearchPageSideBar(props) {
    const {
        onClear,
        onFileUpload
    } = props;

    return (
        <React.Fragment>
            <SideBarCategory>
                <SideBarButton icon={UploadIcon} text="Upload File" onClick={onFileUpload}/>
                <SideBarButton icon={SearchIcon} text="Clear Results" onClick={onClear}/>
            </SideBarCategory>
            <SideBarCategory title="Popular Tags">
                <TagCloud tags={["animals", "birds", "pdf", "image", "novel", "fiction"]} onClick={console.log}/>
            </SideBarCategory>
        </React.Fragment>
    );
}


SearchPageSideBar.propTypes = {
    onClear: PropTypes.func.isRequired,
    onFileUpload: PropTypes.func.isRequired,
};


export default SearchPageSideBar;
