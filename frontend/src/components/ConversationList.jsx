import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import * as conversationApi from '../api/conversation';

const ConversationList = () => {
  const [conversations, setConversations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

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

  return (
      <div className="page-container">
        <div className="page-header">
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

        {error && <div className="error-message">{error}</div>}

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
                  <Link
                      key={conv.id}
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
                  </Link>
              ))}
            </div>
        )}
      </div>
  );
};

export default ConversationList;