import { createContext, useContext, useState, useEffect } from 'react';
import * as authApi from '../api/auth';

const AuthContext = createContext(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    setError(null);
    try {
      const response = await authApi.login({ email, password });
      const userData = {
        id: response.data.id,
        username: response.data.username,
      };
      const token = response.data.accessToken || response.data.access_token;
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
      if (token) {
        localStorage.setItem('access_token', token);
      }
      return userData;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  const register = async (username, email, password) => {
    setError(null);
    try {
      const response = await authApi.register({ username, email, password });
      return response.data;
    } catch (err) {
      setError(err.message);
      throw err;
    }
  };

  const logout = async () => {
    setError(null);
    try {
      await authApi.logout();
      setUser(null);
      localStorage.removeItem('user');
      localStorage.removeItem('access_token');
    } catch (err) {
      setError(err.message);
      setUser(null);
      localStorage.removeItem('user');
      localStorage.removeItem('access_token');
    }
  };

  const value = {
    user,
    loading,
    error,
    login,
    register,
    logout,
    isAuthenticated: !!user,
  };

  return (
      <AuthContext.Provider value={value}>
        {children}
      </AuthContext.Provider>
  );
};
