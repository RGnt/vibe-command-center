import React, { useState } from 'react';
import { useAuth } from '../../contexts/AuthContext';

export default function Login() {
  const { login, register } = useAuth();
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const result = isLogin 
        ? await login(email, password) 
        : await register(email, password);

      if (!result.success) {
        setError(result.error || 'Authentication failed');
      }
    } catch (err) {
      setError('Network error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen w-full bg-[#1a1b26] flex items-center justify-center p-4">
      {/* Dynamic Background Elements */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-[#7aa2f7]/20 rounded-full blur-3xl" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-[#bb9af7]/20 rounded-full blur-3xl" />
      
      <div className="relative w-full max-w-md">
        {/* Glassmorphism Card */}
        <div className="bg-[#1f2335]/80 backdrop-blur-xl rounded-2xl shadow-2xl p-8 border border-[#292e42]">
          
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-[#c0caf5] mb-2 tracking-tight">
              {isLogin ? 'Welcome Back' : 'Create Account'}
            </h1>
            <p className="text-[#565f89]">
              {isLogin ? 'Enter your credentials to continue' : 'Sign up for a new workspace'}
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label className="block text-sm font-medium text-[#9aa5ce] mb-1.5">
                Email Address
              </label>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-[#16161e] border border-[#292e42] rounded-lg px-4 py-3 text-[#c0caf5] placeholder-[#565f89] focus:outline-none focus:ring-2 focus:ring-[#7aa2f7] focus:border-transparent transition-all duration-200"
                placeholder="name@example.com"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-[#9aa5ce] mb-1.5">
                Password
              </label>
              <input
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-[#16161e] border border-[#292e42] rounded-lg px-4 py-3 text-[#c0caf5] placeholder-[#565f89] focus:outline-none focus:ring-2 focus:ring-[#7aa2f7] focus:border-transparent transition-all duration-200"
                placeholder="••••••••"
              />
            </div>

            {error && (
              <div className="bg-[#f7768e]/10 border border-[#f7768e]/20 text-[#f7768e] text-sm rounded-lg p-3 text-center">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-[#7aa2f7] hover:bg-[#7dcfff] text-[#1a1b26] font-semibold py-3 px-4 rounded-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed shadow-[0_0_15px_rgba(122,162,247,0.3)] hover:shadow-[0_0_25px_rgba(122,162,247,0.5)] transform hover:-translate-y-0.5"
            >
              {loading ? 'Processing...' : (isLogin ? 'Sign In' : 'Sign Up')}
            </button>
          </form>

          <div className="mt-6 text-center">
            <button
              onClick={() => { setIsLogin(!isLogin); setError(''); }}
              className="text-sm text-[#bb9af7] hover:text-[#c0caf5] transition-colors duration-200"
            >
              {isLogin ? "Don't have an account? Sign up" : "Already have an account? Sign in"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
