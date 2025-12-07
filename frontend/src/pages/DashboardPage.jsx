import React, { useState, useEffect } from 'react'
import api from '../api'
import jwt_decode from '../jwt-decode'

function DashboardPage() {
  const [tasks, setTasks] = useState([])
  const [sharedTasks, setSharedTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [loadingShared, setLoadingShared] = useState(false)
  const [error, setError] = useState('')
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [deadline, setDeadline] = useState('')
  const [shareUser, setShareUser] = useState('')
  const [shareLoading, setShareLoading] = useState(false)
  const [shareError, setShareError] = useState('')
  const [shareSuccess, setShareSuccess] = useState('')
  const [shareTaskId, setShareTaskId] = useState(null)
  const [userInfo, setUserInfo] = useState({})

  // Carrega as tasks ao montar o componente
  useEffect(() => {
    fetchTasks()
    fetchSharedTasks()
    // Buscar dados do usuário logado
    const fetchUserInfo = async () => {
      try {
        const response = await api.get('/auth/me')
        let username = response.data.username || ''
        if (username.length > 0) {
          username = username.charAt(0).toUpperCase() + username.slice(1)
        }
        let email = response.data.email || ''
        email = email.trim().toLowerCase()
        setUserInfo({
          id: response.data.id,
          username,
          email
        })
      } catch {}
    }
    fetchUserInfo()
  }, [])

  const fetchTasks = async () => {
    try {
      setLoading(true)
      const response = await api.get('/tasks/')
      setTasks(response.data || [])
    } catch (err) {
      setError('Erro ao carregar tasks: ' + (err.response?.data?.error || err.message))
    } finally {
      setLoading(false)
    }
  }

  const fetchSharedTasks = async () => {
    try {
      setLoadingShared(true)
      const response = await api.get('/tasks/shared')
      setSharedTasks(response.data || [])
    } catch (err) {
      if (err.response && err.response.status === 404) {
        setSharedTasks([])
      } else {
        setError('Erro ao carregar tasks compartilhadas: ' + (err.response?.data?.error || err.message))
      }
    } finally {
      setLoadingShared(false)
    }
  }

  const handleCreateTask = async (e) => {
    e.preventDefault()
    if (!title || !deadline) {
      setError('Preencha pelo menos título e deadline')
      return
    }

      let deadlineISO = deadline
      if (deadline) {
        const d = new Date(deadline)
        deadlineISO = d.toISOString()
      }

    try {
      await api.post('/tasks/create', {
        title,
        description,
         deadline: deadlineISO,
        category: 'general'
      })
      setTitle('')
      setDescription('')
      setDeadline('')
      setError('')
      fetchTasks() // Recarrega a lista
    } catch (err) {
      setError('Erro ao criar task: ' + (err.response?.data?.error || err.message))
    }
  }

  const handleCompleteTask = async (id) => {
    try {
      await api.post(`/tasks/complete/${id}`)
      fetchTasks()
    } catch (err) {
      setError('Erro ao completar task: ' + err.message)
    }
  }

  const handleDeleteTask = async (id) => {
    if (window.confirm('Tem certeza?')) {
      try {
        await api.delete(`/tasks/delete/${id}`)
        fetchTasks()
      } catch (err) {
        setError('Erro ao deletar task: ' + err.message)
      }
    }
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    window.location.href = '/login'
  }

  const handleShareTask = async (taskId) => {
    if (!shareUser) {
      setShareError('Informe o usuário ou email para compartilhar')
      return
    }
    setShareLoading(true)
    setShareError('')
    setShareSuccess('')
    try {
      await api.post(`/tasks/share/${taskId}`, { user: shareUser.trim().toLowerCase() })
      setShareSuccess('Tarefa compartilhada com sucesso!')
      setShareUser('')
      setShareTaskId(null)
      fetchTasks()
      fetchSharedTasks()
    } catch (err) {
      setShareError(err.response?.data?.error || 'Erro ao compartilhar tarefa')
    } finally {
      setShareLoading(false)
    }
  }

  return (
    <div style={{minHeight:'100vh',background:'linear-gradient(135deg,#ffe5c2 0%,#ffd8b0 100%)',boxSizing:'border-box',fontFamily:'Poppins, Arial, sans-serif'}}>
      {/* Logo/título no topo */}
      <div style={{width:'100%',padding:'32px 0 0 0',textAlign:'center'}}>
        <span
          style={{
            fontWeight:900,
            fontSize:28,
            letterSpacing:2,
            color:'#ff9800',
            fontFamily:'Poppins, Arial, sans-serif',
            background:'#fff',
            borderRadius:'18px',
            boxShadow:'0 4px 24px #ff980033, 0 1.5px 8px #ff980022',
            padding:'10px 36px',
            display:'inline-block',
            border:'2.5px solid #ffd8b0',
          }}
        >
          To-do API
        </span>
      </div>
      <div style={{display:'flex',flexDirection:'row',gap:32,padding:32,borderRadius:18,marginTop:0}}>
        {/* Coluna esquerda: Minhas tasks */}
        <div style={{flex:1,display:'flex',flexDirection:'column',gap:32}}>
          <div style={{borderRadius:16,boxShadow:'0 2px 12px #0002',background:'#fff',padding:'20px 32px',marginBottom:0,display:'flex',alignItems:'center',justifyContent:'space-between',gap:24}}>
            <span style={{fontWeight:600,fontSize:22,color:'#222'}}>Minhas tarefas</span>
            <div style={{display:'flex',alignItems:'center',gap:16}}>
              <span style={{fontWeight:500,fontSize:16,color:'#555',background:'#fff3e0',padding:'6px 18px',borderRadius:12,boxShadow:'0 1px 4px #ff980011'}}>Usuário logado: <b style={{color:'#ff9800'}}>{userInfo.username || 'Desconhecido'}</b> <span style={{color:'#aaa'}}>(ID: {userInfo.id})</span></span>
              <button onClick={handleLogout} style={{color:'#d32f2f',background:'#fff3e0',border:'none',borderRadius:12,padding:'6px 18px',fontWeight:600,fontSize:16,boxShadow:'0 1px 4px #d32f2f11',cursor:'pointer'}}>Logout</button>
            </div>
          </div>
          {/* Lista de tasks */}
          <div style={{marginTop:0}}>
          {loading ? (
            <p style={{color:'#888'}}>Carregando tasks...</p>
          ) : tasks.length === 0 ? (
            <p style={{color:'#888'}}>Nenhuma task encontrada</p>
          ) : (
            tasks.map((task) => (
              <div key={task.id} style={{borderRadius:16,boxShadow:'0 2px 12px #ff980033',padding:28,minHeight:120,marginBottom:8,background:'#fff',transition:'box-shadow 0.2s',border:'none',display:'flex',flexDirection:'column',gap:8}}>
                <div style={{fontWeight:600,fontSize:18,color:'#222'}}>{task.title}</div>
                <div style={{marginBottom:4,color:'#444'}}>{task.description}</div>
                <div style={{marginBottom:4,fontSize:14,color:'#888'}}>Deadline: {new Date(task.deadline).toLocaleString('pt-BR')}</div>
                {task.shared_with && task.shared_with.length > 0 && (
                  <div style={{marginBottom:4}}>
                    <span style={{fontSize:13,color:'#ff9800'}}>Compartilhada com: {task.shared_with.map((u, i) => (
                      <span key={u}>{u}{i < task.shared_with.length - 1 ? ', ' : ''}</span>
                    ))}</span>
                  </div>
                )}
                <div style={{display:'flex',gap:10,marginTop:8}}>
                  {!task.completed && (
                    <button onClick={() => handleCompleteTask(task.id)} style={{background:'#ff9800',color:'#fff',border:'none',borderRadius:8,padding:'8px 18px',fontWeight:500,boxShadow:'0 1px 4px #ff980033',cursor:'pointer',transition:'background 0.2s'}}>
                      ✓ Completar
                    </button>
                  )}
                  <button onClick={() => handleDeleteTask(task.id)} style={{background:'#d32f2f',color:'#fff',border:'none',borderRadius:8,padding:'8px 18px',fontWeight:500,boxShadow:'0 1px 4px #d32f2f33',cursor:'pointer',transition:'background 0.2s'}}>
                    Deletar
                  </button>
                  <button onClick={() => { setShareTaskId(task.id); setShareUser(''); setShareError(''); setShareSuccess(''); }} style={{background:'#388e3c',color:'#fff',border:'none',borderRadius:8,padding:'8px 18px',fontWeight:500,boxShadow:'0 1px 4px #388e3c33',cursor:'pointer',transition:'background 0.2s'}}>
                    Compartilhar
                  </button>
                </div>
                {shareTaskId === task.id && (
                  <div style={{marginTop:12,display:'flex',gap:8,alignItems:'center'}}>
                    <input
                      type="text"
                      placeholder="Usuário para compartilhar"
                      value={shareUser}
                      onChange={e => setShareUser(e.target.value)}
                      style={{padding:8,borderRadius:8,border:'1px solid #ccc',fontSize:14}}
                    />
                    <button onClick={() => handleShareTask(task.id)} disabled={shareLoading} style={{background:'#388e3c',color:'#fff',border:'none',borderRadius:8,padding:'8px 18px',fontWeight:500,boxShadow:'0 1px 4px #388e3c33',cursor:'pointer'}}>
                      {shareLoading ? 'Compartilhando...' : 'Compartilhar'}
                    </button>
                    {shareError && <span style={{color:'#d32f2f',fontSize:13}}>{shareError}</span>}
                    {shareSuccess && <span style={{color:'#388e3c',fontSize:13}}>{shareSuccess}</span>}
                  </div>
                )}
              </div>
            ))
          )}
          </div>
        </div>
        {/* Coluna direita: Criar task + compartilhadas */}
        <div style={{flex:1,display:'flex',flexDirection:'column',gap:32,justifyContent:'flex-start'}}>
          <div style={{borderRadius:16,boxShadow:'0 2px 12px #ff980033',padding:28,minHeight:180,background:'#fff',marginBottom:0,border:'none',display:'flex',flexDirection:'column',gap:12,marginTop:0}}>
            <div style={{fontWeight:600,fontSize:18,color:'#222',marginBottom:8}}>Criar nova task</div>
            <form onSubmit={handleCreateTask} style={{display:'flex',flexDirection:'column',gap:10}}>
              <input
                type="text"
                placeholder="Título"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                style={{padding:10,borderRadius:8,border:'1px solid #ccc',fontSize:15}}
              />
              <textarea
                placeholder="Descrição"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                style={{padding:10,borderRadius:8,border:'1px solid #ccc',fontSize:15}}
              />
              <input
                type="datetime-local"
                value={deadline}
                onChange={(e) => setDeadline(e.target.value)}
                style={{padding:10,borderRadius:8,border:'1px solid #ccc',fontSize:15}}
              />
              <button type="submit" style={{color:'#fff',background:'#ff9800',border:'none',borderRadius:8,padding:'10px 0',fontWeight:600,fontSize:16,marginTop:8,boxShadow:'0 1px 4px #ff980033',cursor:'pointer'}}>Criar</button>
            </form>
            {error && <p style={{color:'#d32f2f',marginTop:8}}>{error}</p>}
          </div>
          <div style={{borderRadius:16,boxShadow:'0 2px 12px #ff980033',padding:28,minHeight:180,background:'#fff',border:'none',display:'flex',flexDirection:'column',gap:12,marginTop:0}}>
            <div style={{fontWeight:600,fontSize:18,color:'#222',marginBottom:8}}>Tarefas Compartilhadas Comigo</div>
            {loadingShared ? (
              <p style={{color:'#888'}}>Carregando tarefas compartilhadas...</p>
            ) : sharedTasks.length === 0 ? (
              <p style={{color:'#888'}}>Nenhuma tarefa compartilhada encontrada</p>
            ) : (
              sharedTasks.map((task) => (
                <div key={task.id} style={{borderRadius:12,boxShadow:'0 1px 6px #ff980033',padding:18,marginBottom:10,background:'#fff3e0',border:'none',display:'flex',flexDirection:'column',gap:6}}>
                  <div style={{fontWeight:600,fontSize:16,color:'#ff9800'}}>{task.title}</div>
                  <div style={{color:'#444'}}>{task.description}</div>
                  <div style={{fontSize:14,color:'#888'}}>Deadline: {new Date(task.deadline).toLocaleString('pt-BR')}</div>
                  <div style={{fontSize:13,color:'#ff9800'}}>Owner: {task.owner}</div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  )
}

export default DashboardPage
