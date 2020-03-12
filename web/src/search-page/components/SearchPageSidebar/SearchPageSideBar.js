import React from "react";
import SideBarCategory from "./SideBarCategory";
import SideBarButton from "./SideBarButton";
import SearchIcon from "@material-ui/icons/Search";
import UploadIcon from "@material-ui/icons/CloudUpload";
import {useDispatch} from "react-redux";
import {clear} from "../../state";

function SearchPageSideBar() {
    const dispatch = useDispatch();

    return (
        <React.Fragment>
            <SideBarCategory>
                <SideBarButton icon={UploadIcon} text="Upload File"/>
                <SideBarButton icon={SearchIcon} text="Clear Results" onClick={() => dispatch(clear())}/>
            </SideBarCategory>
        </React.Fragment>
    );
}


export default SearchPageSideBar;
