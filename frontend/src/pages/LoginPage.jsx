import React, { useState } from 'react'
import api from '../api'

function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleLogin = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await api.post('/login', {
        username,
        password
      })

      // Salva o token no localStorage
      localStorage.setItem('token', response.data.token)
      
      // Redireciona para dashboard
      window.location.href = '/dashboard'
    } catch (err) {
      setError(err.response?.data?.error || 'Erro ao fazer login')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-container">
      <h1>Login - Todo App</h1>
      <form onSubmit={handleLogin}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit" disabled={loading}>
          {loading ? 'Carregando...' : 'Login'}
        </button>
      </form>
      {error && <p className="error">{error}</p>}
      <p>
        Não tem conta? <a href="/register">Registre-se aqui</a>
      </p>
    </div>
  )
}

export default LoginPage
