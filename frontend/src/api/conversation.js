const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

const getAuthHeaders = () => {
  const token = localStorage.getItem('access_token');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

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

export const createConversation = async (data) => {
  const response = await fetch(`${API_BASE_URL}/conversations`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
    body: JSON.stringify(data),
  });

  return parseResponse(response);
};

export const getConversations = async () => {
  const response = await fetch(`${API_BASE_URL}/conversations`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
  });

  return parseResponse(response);
};

export const getConversation = async (id) => {
  const response = await fetch(`${API_BASE_URL}/conversations/${id}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
  });

  return parseResponse(response);
};

export const deleteConversation = async (id) => {
  const response = await fetch(`${API_BASE_URL}/conversations/${id}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
  });

  return parseResponse(response);
};

export const updateConversation = async (id, data) => {
  const response = await fetch(`${API_BASE_URL}/conversations/${id}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
    body: JSON.stringify(data),
  });

  return parseResponse(response);
};

export const leaveConversation = async (id) => {
  const response = await fetch(`${API_BASE_URL}/conversations/${id}/leave`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
  });

  return parseResponse(response);
};

export const getMessages = async (conversationId, limit = 50, offset = 0) => {
  const response = await fetch(
      `${API_BASE_URL}/conversations/${conversationId}/messages?limit=${limit}&offset=${offset}`,
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders(),
        },
      }
  );

  return parseResponse(response);
};

export const sendMessage = async (conversationId, data) => {
  const response = await fetch(`${API_BASE_URL}/conversations/${conversationId}/messages`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
    body: JSON.stringify(data),
  });

  return parseResponse(response);
};

export const deleteMessage = async (conversationId, messageId) => {
  const response = await fetch(
      `${API_BASE_URL}/conversations/${conversationId}/messages/${messageId}`,
      {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders(),
        },
      }
  );

  return parseResponse(response);
};

export const getAIModels = async () => {
  const response = await fetch(`${API_BASE_URL}/ai/models`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
  });

  return parseResponse(response);
};

export const getAISettings = async (conversationId) => {
  const response = await fetch(`${API_BASE_URL}/conversations/${conversationId}/ai-settings`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
  });

  return parseResponse(response);
};

export const updateAISettings = async (conversationId, data) => {
  const response = await fetch(`${API_BASE_URL}/conversations/${conversationId}/ai-settings`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders(),
    },
    body: JSON.stringify(data),
  });

  return parseResponse(response);
};

export const sendAIMessage = async (conversationId, content, onChunk, onDone, onError) => {
  const token = localStorage.getItem('access_token');
  
  try {
    const response = await fetch(`${API_BASE_URL}/conversations/${conversationId}/ai/stream`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
      },
      body: JSON.stringify({ content }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `Request failed with status ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop() || '';

      for (const line of lines) {
        const trimmed = line.trim();
        if (!trimmed || !trimmed.startsWith('data: ')) continue;
        
        const data = trimmed.slice(6);
        if (data === '[DONE]') {
          if (onDone) onDone();
          return;
        }

        if (onChunk) onChunk(data);
      }
    }

    if (onDone) onDone();
  } catch (err) {
    if (onError) onError(err.message);
    else throw err;
  }
};