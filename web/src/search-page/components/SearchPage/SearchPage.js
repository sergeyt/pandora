import React from "react";
import Dashboard from "../Dashboard";
import SearchPageSideBar from "../SearchPageSidebar";
import QueryInput from "../QueryInput";


function SearchPage() {
    return (
        <Dashboard
            sidebar={
                <SearchPageSideBar/>
            }
            content={
                <QueryInput/>
            }
        />
    );
}

export default SearchPage;
