import React from "react";
import Dashboard from "../Dashboard";
import SearchPageSideBar from "../SearchPageSidebar";
import QueryInput from "../QueryInput";
import {useSelector} from "react-redux";
import SearchResults from "../SearchResults";


function SearchPage() {
    const documents = useSelector(state => state.search.documents);

    const content = (documents.length == 0 ? <QueryInput/> : <SearchResults/>);

    return (
        <Dashboard
            sidebar={<SearchPageSideBar/>}
            content={content}
        />
    );
}

export default SearchPage;
