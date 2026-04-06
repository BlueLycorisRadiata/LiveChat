import { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import * as conversationApi from '../api/conversation';

const CreateConversation = () => {
  const [title, setTitle] = useState('');
  const [type, setType] = useState('private');
  const [model, setModel] = useState('');
  const [models, setModels] = useState([]);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    const fetchModels = async () => {
      try {
        const response = await conversationApi.getAIModels();
        setModels(response.data || []);
        if (response.data && response.data.length > 0) {
          setModel(response.data[0].id);
        }
      } catch (err) {
        console.error('Failed to fetch models:', err);
      }
    };

    fetchModels();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess(false);

    setIsLoading(true);
    try {
      const response = await conversationApi.createConversation({
        title: title.trim() || null,
        type,
        model: type === 'ai' ? model : null,
        participant_ids: [],
      });
      setSuccess(true);
      setTimeout(() => {
        navigate(`/conversations/${response.data.id}`);
      }, 1500);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
      <div className="page-container">
        <div className="page-header">
          <Link to="/conversations" className="back-link">
            &larr; Back to Conversations
          </Link>
          <h1>New Conversation</h1>
        </div>

        <div className="form-card">
          <h2>New Conversation</h2>
          {error && <div className="error-message">{error}</div>}
          {success && (
              <div className="success-message">
                Conversation created! Redirecting...
              </div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="title">Conversation Name (optional)</label>
              <input
                  type="text"
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="Enter conversation name"
              />
            </div>

            <div className="form-group">
              <label htmlFor="type">Conversation Type</label>
              <select
                  id="type"
                  value={type}
                  onChange={(e) => setType(e.target.value)}
              >
                <option value="private">Private (1-on-1)</option>
                <option value="group">Group</option>
                <option value="ai">AI Chat</option>
              </select>
            </div>

            {type === 'ai' && (
                <div className="form-group">
                  <label htmlFor="model">AI Model</label>
                  <select
                      id="model"
                      value={model}
                      onChange={(e) => setModel(e.target.value)}
                  >
                    {models.map((m) => (
                        <option key={m.id} value={m.id}>
                          {m.name}
                        </option>
                    ))}
                  </select>
                </div>
            )}

            <button type="submit" disabled={isLoading} className="submit-button">
              {isLoading ? 'Creating...' : 'Create Conversation'}
            </button>
          </form>
        </div>
      </div>
  );
};

export default CreateConversation;