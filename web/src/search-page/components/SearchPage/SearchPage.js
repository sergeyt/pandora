import React from "react";
import Dashboard from "../Dashboard";
import SearchPageSideBar from "../SearchPageSidebar";


function SearchPage() {
    return (
        <Dashboard
            sidebar={
                <SearchPageSideBar/>
            }
            content={
                <h1>Hello Pandora</h1>
            }
        />
    );
}

export default SearchPage;
