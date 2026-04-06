const parseResponse = async (response) => {
  const contentType = response.headers.get('content-type');
  const isJson = contentType && contentType.includes('application/json');

  const payload = isJson ? await response.json() : null;

  if (!response.ok) {
    const message =
        payload?.error ||
        payload?.message ||
        `Request failed with status ${response.status}`;
    throw new Error(message);
  }

  return payload;
};

const getAuthHeaders = () => {
  const token = localStorage.getItem('access_token');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

export const register = async (userData) => {
  const response = await fetch('/signup', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  return parseResponse(response);
};

export const login = async (credentials) => {
  const response = await fetch('/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(credentials),
  });

  return parseResponse(response);
};

export const logout = async () => {
  const token = localStorage.getItem('access_token');
  const headers = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  try {
    const response = await fetch('/logout', {
      method: 'POST',
      headers,
    });
    // Don't need to parse response - just clear local storage
  } catch (e) {
    // Ignore errors
  } finally {
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
  }
};