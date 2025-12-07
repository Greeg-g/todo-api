import React, { useState } from 'react'
import api from '../api'

function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleLogin = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await api.post('/auth/login', {
        email,
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
    <div className="login-container" style={{maxWidth:400,margin:'40px auto',background:'#fff',borderRadius:8,boxShadow:'0 2px 8px #0001',padding:32}}>
      <h1 style={{textAlign:'center',marginBottom:24}}>Login - Todo App</h1>
      <form onSubmit={handleLogin} style={{display:'flex',flexDirection:'column',gap:16}}>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          required
          style={{padding:10,borderRadius:4,border:'1px solid #ccc'}}
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          style={{padding:10,borderRadius:4,border:'1px solid #ccc'}}
        />
        <button type="submit" disabled={loading} style={{padding:10,borderRadius:4,background:'#1976d2',color:'#fff',border:'none',fontWeight:'bold'}}>
          {loading ? 'Carregando...' : 'Login'}
        </button>
      </form>
      {error && <p className="error" style={{color:'red',marginTop:12}}>{error}</p>}
      <p style={{marginTop:24,textAlign:'center'}}>
        Não tem conta? <a href="/register" style={{color:'#1976d2'}}>Registre-se aqui</a>
      </p>
    </div>
  )
}

export default LoginPage
