const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

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
  const response = await fetch(`${API_BASE_URL}/signup`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });

  return parseResponse(response);
};

export const login = async (credentials) => {
  const response = await fetch(`${API_BASE_URL}/login`, {
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
    const response = await fetch(`${API_BASE_URL}/logout`, {
      method: 'POST',
      headers,
    });
  } catch (e) {
  } finally {
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
  }
};