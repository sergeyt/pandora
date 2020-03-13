export const ACTION_QUERY = "search.QUERY";
export const ACTION_QUERY_RESULTS = "search.QUERY_RESULTS";
export const ACTION_CLEAR = "search.CLEAR";

export function query(queryString) {
    return {
        type: ACTION_QUERY,
        query: queryString,
    };
}

export function queryResults(query, documents, success = true) {
    return {
        type: ACTION_QUERY_RESULTS,
        query,
        documents,
        success,
    };
}

export function clear() {
    return {type: ACTION_CLEAR};
}