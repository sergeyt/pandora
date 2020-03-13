import React, {useState} from "react";
import Dashboard from "../Dashboard";
import SearchPageSideBar from "../SearchPageSidebar";
import QueryInput from "../QueryInput";
import {useDispatch, useSelector} from "react-redux";
import SearchResults from "../SearchResults";
import {clear} from "../../state";
import UploadDialog from "../UploadDialog";


function SearchPage() {
    const documents = useSelector(state => state.search.documents);
    const dispatch = useDispatch();

    // Upload dialog state.
    const [uploadDialogOpen, setUploadDialogOpen] = useState(false);

    // Emphasize search form when results are empty
    const content = (documents.length === 0 ? <QueryInput/> : <SearchResults/>);

    return (
        <React.Fragment>
            <Dashboard
                sidebar={
                    <SearchPageSideBar
                        onClear={() => dispatch(clear())}
                        onFileUpload={() => setUploadDialogOpen(true)}
                    />
                }
                content={content}
            />
            <UploadDialog
                open={uploadDialogOpen}
                onClose={() => setUploadDialogOpen(false)}
                onFileUpload={() => setUploadDialogOpen(false)}
            />
        </React.Fragment>
    );
}

export default SearchPage;
