import React, { useState } from 'react'
import api from '../api'

function RegisterPage() {
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const [loading, setLoading] = useState(false)

  const handleRegister = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')
    setLoading(true)

    try {
      await api.post('/auth/register', {
        username,
        email,
        password
      })

      setSuccess('Conta criada com sucesso! Redirecionando para login...')
      setTimeout(() => {
        window.location.href = '/login'
      }, 2000)
    } catch (err) {
      setError(err.response?.data?.error || 'Erro ao registrar')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="register-container" style={{maxWidth:400,margin:'40px auto',background:'#fff',borderRadius:8,boxShadow:'0 2px 8px #0001',padding:32}}>
      <h1 style={{textAlign:'center',marginBottom:24}}>Registrar - Todo App</h1>
      <form onSubmit={handleRegister} style={{display:'flex',flexDirection:'column',gap:16}}>
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          style={{padding:10,borderRadius:4,border:'1px solid #ccc'}}
        />
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
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
        <button type="submit" disabled={loading} style={{padding:10,borderRadius:4,background:'#388e3c',color:'#fff',border:'none',fontWeight:'bold'}}>
          {loading ? 'Carregando...' : 'Registrar'}
        </button>
      </form>
      {error && <p className="error" style={{color:'red',marginTop:12}}>{error}</p>}
      {success && <p className="success" style={{color:'green',marginTop:12}}>{success}</p>}
      <p style={{marginTop:24,textAlign:'center'}}>
        Já tem conta? <a href="/login" style={{color:'#388e3c'}}>Faça login</a>
      </p>
    </div>
  )
}

export default RegisterPage
