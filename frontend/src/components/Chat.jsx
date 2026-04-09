import { useState, useEffect, useRef, useCallback } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import * as conversationApi from '../api/conversation';

const Chat = () => {
  const { conversationId } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  const [conversation, setConversation] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [sending, setSending] = useState(false);
  const [aiSettings, setAiSettings] = useState(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editTitle, setEditTitle] = useState('');

  const messagesEndRef = useRef(null);

  const fetchConversation = useCallback(async () => {
    try {
      const convRes = await conversationApi.getConversation(conversationId);
      setConversation(convRes.data);
    } catch (err) {
      console.error('Failed to fetch conversation:', err);
    }
  }, [conversationId]);

  const fetchMessages = useCallback(async () => {
    try {
      const msgRes = await conversationApi.getMessages(conversationId);
      return msgRes.data || [];
    } catch (err) {
      console.error('Failed to fetch messages:', err);
      return null;
    }
  }, [conversationId]);

  const fetchAISettings = useCallback(async () => {
    try {
      const res = await conversationApi.getAISettings(conversationId);
      setAiSettings(res.data);
    } catch (err) {
      console.error('Failed to fetch AI settings:', err);
    }
  }, [conversationId]);

  useEffect(() => {
    const loadInitialData = async () => {
      setLoading(true);
      try {
        const [convRes, msgRes] = await Promise.all([
          conversationApi.getConversation(conversationId),
          conversationApi.getMessages(conversationId),
        ]);
        setConversation(convRes.data);
        setMessages(msgRes.data || []);
        setError('');
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    loadInitialData();
  }, [conversationId]);

  useEffect(() => {
    if (conversation?.type === 'ai') {
      fetchAISettings();
    }
  }, [conversation?.type === 'ai', fetchAISettings]);

  useEffect(() => {
    const interval = setInterval(async () => {
      try {
        await conversationApi.getConversation(conversationId);
        const newMessages = await fetchMessages();
        if (newMessages !== null) {
          setMessages(newMessages);
        }
      } catch (err) {
        if (err.message.includes('not found') || err.message.includes('404')) {
          navigate('/conversations');
        }
      }
    }, 3000);
    return () => clearInterval(interval);
  }, [conversationId, fetchMessages, navigate]);
  //
  // useEffect(() => {
  //   messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  // }, [messages]);

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    setSending(true);
    try {
      if (conversation?.type === 'ai') {
        const userMessage = newMessage.trim();
        setNewMessage('');
        
        await conversationApi.sendAIMessage(
          conversationId,
          userMessage,
          (chunk) => {
            setMessages(prev => {
              const lastMsg = prev[prev.length - 1];
              if (lastMsg && lastMsg.role === 'assistant') {
                return [...prev.slice(0, -1), { ...lastMsg, content: lastMsg.content + chunk }];
              }
              return [...prev, {
                id: Date.now(),
                content: chunk,
                role: 'assistant',
                created_at: new Date().toISOString(),
              }];
            });
          },
          () => {},
          (error) => {
            setError(error);
          }
        );
      } else {
        await conversationApi.sendMessage(conversationId, {
          content: newMessage.trim(),
          type: 'text',
        });
        setNewMessage('');
        const newMessages = await fetchMessages();
        if (newMessages !== null) {
          setMessages(newMessages);
        }
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setSending(false);
    }
  };

  const handleDeleteConversation = async () => {
    if (!window.confirm('Are you sure you want to delete this conversation?')) {
      return;
    }

    try {
      await conversationApi.deleteConversation(conversationId);
      navigate('/conversations');
    } catch (err) {
      setError(err.message);
    }
  };

  const handleLeaveConversation = async () => {
    if (!window.confirm('Are you sure you want to leave this conversation?')) {
      return;
    }

    try {
      await conversationApi.leaveConversation(conversationId);
      navigate('/conversations');
    } catch (err) {
      setError(err.message);
    }
  };

  const handleModelChange = async (newModel) => {
    try {
      await conversationApi.updateAISettings(conversationId, { model: newModel });
      setAiSettings(prev => ({ ...prev, model: newModel }));
    } catch (err) {
      setError(err.message);
    }
  };

  const handleEditTitle = async () => {
    if (!editTitle.trim()) return;
    try {
      const res = await conversationApi.updateConversation(conversationId, { title: editTitle.trim() });
      setConversation(prev => ({ ...prev, title: res.data.title }));
      setIsEditing(false);
    } catch (err) {
      setError(err.message);
    }
  };

  const isAIConversation = conversation?.type === 'ai';

  const getMessageClass = (msg) => {
    if (isAIConversation) {
      return msg.role === 'assistant' ? 'ai' : 'user';
    }
    return msg.sender_id === parseInt(user?.id) ? 'own' : 'other';
  };

  const getMessageSender = (msg) => {
    if (isAIConversation) {
      return msg.role === 'assistant' ? 'AI Assistant' : 'You';
    }
    return msg.sender?.username || 'Unknown';
  };

  if (loading) {
    return (
      <div className="ai-chat-container">
        <div className="ai-chat-header">
          <Link to="/" className="back-link">
            ← Home
          </Link>
        </div>
        <div className="ai-empty-state">
          <div className="loading">Loading conversation...</div>
        </div>
      </div>
    );
  }

  if (isAIConversation) {
    return (
      <div className="ai-chat-container">
        <div className="ai-chat-header">
          <Link to="/" className="back-link">
            ← Home
          </Link>
          <h1>AI Chat</h1>
          {aiSettings && (
            <div className="ai-controls">
              <select 
                value={aiSettings.model} 
                onChange={(e) => handleModelChange(e.target.value)}
                className="model-select"
              >
                <option value="nvidia/nemotron-nano-9b-v2:free">Nemotron Nano</option>
                <option value="minimax/minimax-m2.5:free">Minimax M2.5</option>
                <option value="qwen/qwen3.6-plus:free">Qwen 3.6+</option>
                <option value="nvidia/nemotron-3-super-120b-a12b:free">Nemotron Super</option>
              </select>
              <button onClick={handleDeleteConversation} className="delete-button">
                Delete
              </button>
            </div>
          )}
        </div>

        {error && <div className="error-message" style={{ margin: '16px 24px', background: 'rgba(220,53,69,0.2)', color: '#ff6b6b' }}>{error}</div>}

        <div className="ai-messages-wrapper">
          <div className="ai-messages">
            {messages.length === 0 ? (
              <div className="ai-empty-state">
                <div className="robot-icon">🤖</div>
                <p>Start chatting with AI!</p>
              </div>
            ) : (
              messages.map((msg) => (
                <div key={msg.id} className={`ai-message ${msg.role === 'assistant' ? 'assistant' : 'user'}`}>
                  <div className="ai-message-avatar">
                    {msg.role === 'assistant' ? '🤖' : '👤'}
                  </div>
                  <div>
                    <div className="ai-message-content">{msg.content}</div>
                    <div className="ai-message-time">
                      {new Date(msg.created_at).toLocaleTimeString()}
                    </div>
                  </div>
                </div>
              ))
            )}
            <div ref={messagesEndRef} />
          </div>
        </div>

        <div className="ai-input-wrapper">
          <form onSubmit={handleSendMessage} className="ai-input-container">
            <input
              type="text"
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              placeholder="Ask AI anything..."
              disabled={sending}
            />
            <button type="submit" disabled={sending || !newMessage.trim()}>
              ➤
            </button>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <div className="page-header">
        <Link to="/" className="back-link">
          ← Home
        </Link>
        {isEditing ? (
          <div className="edit-title-form">
            <input
              type="text"
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              placeholder="Enter new title"
              autoFocus
              onKeyDown={(e) => e.key === 'Enter' && handleEditTitle()}
            />
            <button onClick={handleEditTitle} className="submit-button">Save</button>
            <button onClick={() => setIsEditing(false)} className="cancel-button">Cancel</button>
          </div>
        ) : (
          <h1 onClick={() => { setEditTitle(conversation?.title || ''); setIsEditing(true); }} style={{ cursor: 'pointer' }}>
            {conversation?.title || `Conversation ${conversationId}`}
          </h1>
        )}
        <div className="header-actions">
          <Link to="/conversations" className="back-link" style={{ marginRight: '12px' }}>
            ← List
          </Link>
          <button onClick={handleLeaveConversation} className="leave-button">
            Leave
          </button>
          <button onClick={handleDeleteConversation} className="delete-button">
            Delete
          </button>
        </div>
      </div>

      {error && <div className="error-message">{error}</div>}

      <div className="chat-container">
        <div className="chat-messages">
          {messages.length === 0 ? (
            <div className="empty-chat">
              <p>No messages yet. Start the conversation!</p>
            </div>
          ) : (
            messages.map((msg) => (
              <div key={msg.id} className={`message ${getMessageClass(msg)}`}>
                <div className="message-header">
                  <span className="message-username">{getMessageSender(msg)}</span>
                  <span className="message-time">
                    {new Date(msg.created_at).toLocaleTimeString()}
                  </span>
                </div>
                <div className="message-content">{msg.content}</div>
              </div>
            ))
          )}
          <div ref={messagesEndRef} />
        </div>

        <form onSubmit={handleSendMessage} className="chat-input-form">
          <input
            type="text"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Type a message..."
            disabled={sending}
          />
          <button type="submit" disabled={sending || !newMessage.trim()}>
            {sending ? 'Sending...' : 'Send'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default Chat;