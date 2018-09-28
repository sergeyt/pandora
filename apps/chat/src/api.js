const stdHeaders = {
    'Accept': 'application/json',
    'Content-Type': 'application/json',
};

function getToken() {
    return localStorage.getItem('acess_token');
}

function setToken(value) {
    if (value) {
        localStorage.setItem('access_token', value);
    } else {
        localStorage.removeItem('access_token');
    }
}

function fetchJSON(url, payload, options = {}) {
    const token = getToken();

    return fetch(url, {
        method: options.method || (payload !== undefined ? 'POST' : 'GET'),
        headers: {
            ...stdHeaders,
            'Authorization': token ? `Bearer ${token}` : undefined,
        },
        body: payload !== undefined ? JSON.stringify(payload) : undefined,
        ...options,
    }).then(resp => resp.json());
}

function post(url, payload, options = {}) {
    return fetchJSON(url, payload, {
        ...options,
        method: 'POST',
    });
}

function login(username, password) {
    return post('/api/login', { body: {username, password} }).then(resp => {
        setToken(resp.access_token);
        return resp;
    });
}
