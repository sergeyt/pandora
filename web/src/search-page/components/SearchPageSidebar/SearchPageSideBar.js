import React from "react";
import SideBarCategory from "./SideBarCategory";
import SideBarButton from "./SideBarButton";
import SearchIcon from "@material-ui/icons/Search"
import UploadIcon from "@material-ui/icons/CloudUpload";

function SearchPageSideBar() {
    return (
        <React.Fragment>
            <SideBarCategory>
                <SideBarButton icon={SearchIcon} text="Search"/>
                <SideBarButton icon={UploadIcon} text="Upload File"/>
            </SideBarCategory>
        </React.Fragment>
    );
}


export default SearchPageSideBar;
