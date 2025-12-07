# Frontend - Todo App

Frontend em React + Vite para o projeto Todo API.

## Primeiros Passos

### 1. Instalar dependências
```bash
cd frontend
npm install
```

### 2. Iniciar o servidor de desenvolvimento
```bash
npm run dev
```

O frontend estará disponível em `http://localhost:3000`

## Estrutura do Projeto

```
frontend/
├── src/
│   ├── api.js              # Configuração do axios para chamadas à API
│   ├── App.jsx             # Componente principal com rotas
│   ├── main.jsx            # Entry point do React
│   ├── index.css           # Estilos globais
│   └── pages/
│       ├── LoginPage.jsx   # Página de login
│       ├── RegisterPage.jsx # Página de registro
│       └── DashboardPage.jsx # Página principal de tasks
├── index.html
├── vite.config.js
└── package.json
```

## O que cada página faz

### LoginPage
- Faz login com username e password
- Salva o token JWT no localStorage
- Redireciona para dashboard após sucesso

### RegisterPage
- Cria uma nova conta com username, email e password
- Redireciona para login após sucesso

### DashboardPage
- Lista todas as tasks do usuário
- Cria novas tasks
- Marca tasks como completadas
- Deleta tasks

## Como funciona a comunicação com o backend

O arquivo `src/api.js` configura:
- URL base da API: `http://localhost:8080`
- Adiciona o token JWT automaticamente em todas as requisições
- Busca o token do localStorage

Exemplo de como fazer requisições:
```javascript
import api from './api'

// GET
const response = await api.get('/tasks')

// POST
await api.post('/tasks/create', { title, description, deadline })

// DELETE
await api.delete('/tasks/delete/1')
```

## Para rodar tudo junto

### Terminal 1 - Backend
```bash
cd backend
go run main.go
```

### Terminal 2 - Frontend
```bash
cd frontend
npm install (primeira vez)
npm run dev
```

Agora você tem:
- Backend: `http://localhost:8080`
- Frontend: `http://localhost:3000`
