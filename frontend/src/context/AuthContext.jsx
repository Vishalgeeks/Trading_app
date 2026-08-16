import { createContext, useContext, useState, useEffect } from 'react';
import { authService } from '../services/authService';

const AuthContext = createContext(null);

const TOKEN_KEY = 'token';
const USER_KEY = 'user';

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const initializeAuth = async () => {
    const token = localStorage.getItem(TOKEN_KEY);
    if (!token) {
      setLoading(false);
      return;
    }

    const result = await authService.getCurrentUser();
    if (result.error) {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      setLoading(false);
      return;
    }

    const userData = result.data;
    setUser(userData);
    setIsAuthenticated(true);
    localStorage.setItem(USER_KEY, JSON.stringify(userData));
    setLoading(false);
  };

  useEffect(() => {
    initializeAuth();
  }, []);

  const login = async (email, password) => {
    const result = await authService.login(email, password);
    if (result.error) {
      return { error: true, message: result.message };
    }

    const { access_token, user: userData } = result.data;
    localStorage.setItem(TOKEN_KEY, access_token);
    localStorage.setItem(USER_KEY, JSON.stringify(userData));
    setUser(userData);
    setIsAuthenticated(true);

    return { error: false, user: userData };
  };

  const register = async (name, email, phone, password) => {
    const result = await authService.register(name, email, phone, password);
    if (result.error) {
      return { error: true, message: result.message };
    }

    const { access_token, user: userData } = result.data;
    localStorage.setItem(TOKEN_KEY, access_token);
    localStorage.setItem(USER_KEY, JSON.stringify(userData));
    setUser(userData);
    setIsAuthenticated(true);

    return { error: false, user: userData };
  };

  const logout = async () => {
    await authService.logout();
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    setUser(null);
    setIsAuthenticated(false);
  };

  const value = {
    user,
    loading,
    isAuthenticated,
    login,
    register,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
