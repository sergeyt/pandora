export const ACTION_QUERY = "search.QUERY";
export const ACTION_QUERY_RESULTS = "search.QUERY_RESULTS";

export function query(queryString) {
    return {
        type: ACTION_QUERY,
        query: queryString,
    };
}

export function queryResults(documents, success = true) {
    return {
        type: ACTION_QUERY_RESULTS,
        documents,
        success,
    };
}
