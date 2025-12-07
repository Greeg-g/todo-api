import React, { useState, useEffect } from 'react'
import api from '../api'

function DashboardPage() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [deadline, setDeadline] = useState('')

  // Carrega as tasks ao montar o componente
  useEffect(() => {
    fetchTasks()
  }, [])

  const fetchTasks = async () => {
    try {
      setLoading(true)
      const response = await api.get('/tasks')
      setTasks(response.data || [])
    } catch (err) {
      setError('Erro ao carregar tasks: ' + (err.response?.data?.error || err.message))
    } finally {
      setLoading(false)
    }
  }

  const handleCreateTask = async (e) => {
    e.preventDefault()
    if (!title || !deadline) {
      setError('Preencha pelo menos título e deadline')
      return
    }

    try {
      await api.post('/tasks/create', {
        title,
        description,
        deadline,
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

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>Minhas Tasks</h1>
        <button onClick={handleLogout} className="logout-btn">Logout</button>
      </div>

      <form onSubmit={handleCreateTask} className="create-task-form">
        <h2>Criar Nova Task</h2>
        <input
          type="text"
          placeholder="Título"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
        <textarea
          placeholder="Descrição"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
        <input
          type="datetime-local"
          value={deadline}
          onChange={(e) => setDeadline(e.target.value)}
        />
        <button type="submit">Criar Task</button>
      </form>

      {error && <p className="error">{error}</p>}

      {loading ? (
        <p>Carregando tasks...</p>
      ) : tasks.length === 0 ? (
        <p>Nenhuma task encontrada</p>
      ) : (
        <div className="tasks-list">
          {tasks.map((task) => (
            <div key={task.id} className={`task-item ${task.completed ? 'completed' : ''}`}>
              <div className="task-info">
                <h3>{task.title}</h3>
                <p>{task.description}</p>
                <small>Deadline: {new Date(task.deadline).toLocaleString('pt-BR')}</small>
              </div>
              <div className="task-actions">
                {!task.completed && (
                  <button onClick={() => handleCompleteTask(task.id)} className="complete-btn">
                    ✓ Completar
                  </button>
                )}
                <button onClick={() => handleDeleteTask(task.id)} className="delete-btn">
                  Deletar
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default DashboardPage
