import React, {useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {upload, UploadDialog, UploadReport} from "../../../file-upload";
import Dashboard from "../Dashboard";
import SearchPageSideBar from "../SearchPageSidebar";
import QueryInput from "../QueryInput";
import SearchResults from "../SearchResults";
import {clear} from "../../state";

function SearchPage() {
    const documents = useSelector(state => state.search.documents);
    const dispatch = useDispatch();

    // Dashboard sidebar state:
    const [showSidebar, setShowSidebar] = useState(false);

    // Upload dialog state.
    const [uploadDialogOpen, setUploadDialogOpen] = useState(false);

    // Emphasize search form when results are empty
    const content = (documents.length === 0 ? <QueryInput/> : <SearchResults/>);

    const handleUpload = files => {
        setUploadDialogOpen(false);
        dispatch(upload(files));
    };

    const handleShowUploadDialog = () => {
        setShowSidebar(false);
        setUploadDialogOpen(true);
    };

    const handleClear = () => {
        setShowSidebar(false);
        dispatch(clear());
    };

    const handleQuery = (query) => {
        setShowSidebar(false);
        console.log(query);
    };

    return (
        <React.Fragment>
            <Dashboard
                sidebar={
                    <SearchPageSideBar
                        onClear={handleClear}
                        onFileUpload={handleShowUploadDialog}
                        onQuery={handleQuery}
                    />
                }
                content={content}
                showSidebar={showSidebar}
                onShowSidebar={setShowSidebar}
            />
            <UploadDialog
                open={uploadDialogOpen}
                onClose={() => setUploadDialogOpen(false)}
                onFileUpload={handleUpload}
            />
            <UploadReport/>
        </React.Fragment>
    );
}

export default SearchPage;
