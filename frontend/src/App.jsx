import React from 'react';
import TodoApp from './components/TodoApp';
import Login from './components/auth/Login';
import { AuthProvider, useAuth } from './contexts/AuthContext';

function AppContent() {
  const { user } = useAuth();

  if (!user) {
    return <Login />;
  }

  return <TodoApp />;
}

function App() {
  return (
    <div className="min-h-screen bg-bg-base transition-colors duration-300">
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </div>
  );
}

export default App;
