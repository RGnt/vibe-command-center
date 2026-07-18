import React, { createContext, useContext, useState, useEffect } from 'react';

// Patch window.fetch globally to intercept and add the JWT token to all requests
const originalFetch = window.fetch;
window.fetch = async (url, options = {}) => {
  const token = localStorage.getItem('token');
  if (token && (url.startsWith('/api/') || url.startsWith('http'))) {
    options.headers = {
      ...options.headers,
      'Authorization': `Bearer ${token}`
    };
  }
  
  const response = await originalFetch(url, options);
  
  // If unauthorized and we're not hitting auth routes, log out automatically
  if (response.status === 401 && !url.includes('/api/auth/')) {
    localStorage.removeItem('token');
    window.dispatchEvent(new Event('auth:unauthorized'));
  }
  
  return response;
};

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const fetchMe = async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      setLoading(false);
      return;
    }

    try {
      const res = await fetch('/api/auth/me');
      if (res.ok) {
        const data = await res.json();
        setUser(data);
      } else {
        localStorage.removeItem('token');
        setUser(null);
      }
    } catch (e) {
      console.error("Failed to fetch user", e);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMe();
    
    const handleUnauthorized = () => {
      setUser(null);
    };
    window.addEventListener('auth:unauthorized', handleUnauthorized);
    
    return () => window.removeEventListener('auth:unauthorized', handleUnauthorized);
  }, []);

  const login = async (email, password) => {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    
    if (res.ok) {
      const data = await res.json();
      localStorage.setItem('token', data.token);
      setUser(data.user);
      return { success: true };
    }
    
    const err = await res.text();
    return { success: false, error: err };
  };

  const register = async (email, password) => {
    const res = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    });
    
    if (res.ok) {
      const data = await res.json();
      localStorage.setItem('token', data.token);
      setUser(data.user);
      return { success: true };
    }
    
    const err = await res.text();
    return { success: false, error: err };
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
  };

  if (loading) {
    return <div className="h-screen w-screen bg-[#1a1b26] flex items-center justify-center text-[#c0caf5]">Loading...</div>;
  }

  return (
    <AuthContext.Provider value={{ user, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
