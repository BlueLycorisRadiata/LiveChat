import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import * as conversationApi from '../api/conversation';

const ConversationList = () => {
  const [conversations, setConversations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editTitle, setEditTitle] = useState('');
  const navigate = useNavigate();

  const fetchConversations = async () => {
    try {
      setLoading(true);
      const response = await conversationApi.getConversations();
      setConversations(response.data || []);
      setError('');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConversations();
  }, []);

  const handleDelete = async (e, convId) => {
    e.preventDefault();
    e.stopPropagation();
    if (!window.confirm('Are you sure you want to delete this conversation?')) {
      return;
    }
    try {
      await conversationApi.deleteConversation(convId);
      setConversations(prev => prev.filter(c => c.id !== convId));
    } catch (err) {
      setError(err.message);
    }
  };

  const handleEdit = (e, conv) => {
    e.preventDefault();
    e.stopPropagation();
    setEditingId(conv.id);
    setEditTitle(conv.title || '');
  };

  const handleSaveEdit = async (convId) => {
    if (!editTitle.trim()) return;
    try {
      const res = await conversationApi.updateConversation(convId, { title: editTitle.trim() });
      setConversations(prev => prev.map(c => c.id === convId ? { ...c, title: res.data.title } : c));
      setEditingId(null);
    } catch (err) {
      setError(err.message);
    }
  };

  const handleCancelEdit = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setEditingId(null);
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <Link to="/" className="back-link">
          ← Home
        </Link>
        <h1>Conversations</h1>
        <div className="header-actions">
          <button onClick={fetchConversations} className="refresh-button">
            Refresh
          </button>
          <Link to="/conversations/create" className="create-button">
            New Conversation
          </Link>
        </div>
      </div>

      {error && <div className="error-message" style={{ margin: '20px 40px' }}>{error}</div>}

      {loading ? (
        <div className="loading">Loading conversations...</div>
      ) : conversations.length === 0 ? (
        <div className="empty-state">
          <p>No conversations yet. Start a new one!</p>
          <Link to="/conversations/create" className="create-button">
            New Conversation
          </Link>
        </div>
      ) : (
        <div className="conversation-list">
          {conversations.map((conv) => (
            <div key={conv.id} className="conversation-card-wrapper">
              {editingId === conv.id ? (
                <div className="edit-title-form" onClick={(e) => e.stopPropagation()}>
                  <input
                    type="text"
                    value={editTitle}
                    onChange={(e) => setEditTitle(e.target.value)}
                    placeholder="Enter new title"
                    onKeyDown={(e) => e.key === 'Enter' && handleSaveEdit(conv.id)}
                    autoFocus
                  />
                  <button onClick={() => handleSaveEdit(conv.id)} className="submit-button">Save</button>
                  <button onClick={handleCancelEdit} className="cancel-button">Cancel</button>
                </div>
              ) : (
                <Link
                  to={`/conversations/${conv.id}`}
                  className="conversation-card"
                >
                  <div className={`conversation-icon ${conv.type}`}>
                    {conv.type === 'ai' ? '🤖' : conv.type === 'group' ? '👥' : '💬'}
                  </div>
                  <div className="conversation-info">
                    <h3>{conv.title || `Conversation ${conv.id}`}</h3>
                    <p className="conversation-type">{conv.type}</p>
                  </div>
                  <div className="conversation-time">
                    {new Date(conv.updated_at).toLocaleDateString()}
                  </div>
                  <div className="conversation-actions">
                    <button onClick={(e) => handleEdit(e, conv)} className="edit-button" title="Edit">
                      ✏️
                    </button>
                    <button onClick={(e) => handleDelete(e, conv.id)} className="delete-conv-button" title="Delete">
                      🗑️
                    </button>
                  </div>
                </Link>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ConversationList;