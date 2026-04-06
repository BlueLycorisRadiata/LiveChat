import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Home = () => {
  const { user, logout } = useAuth();

  const handleLogout = async () => {
    await logout();
  };

  return (
    <div className="home-container">
      <header className="home-header">
        <h1>Welcome to LiveChat</h1>
        <button onClick={handleLogout} className="logout-button">
          Logout
        </button>
      </header>
      
      <main className="home-content">
        <div className="user-info">
          <h2>Hello, {user?.username}!</h2>
          <p>You are now logged in.</p>
        </div>
        
        <div className="action-cards">
          <Link to="/conversations" className="action-card">
            <h3>Conversations</h3>
            <p>View and manage your conversations</p>
          </Link>
          <Link to="/conversations/create" className="action-card">
            <h3>New Conversation</h3>
            <p>Start a new conversation</p>
          </Link>
        </div>
      </main>
    </div>
  );
};

export default Home;
