# API de Tarefas (To-Do List) em Go

Este projeto é uma API REST desenvolvida em Go com o objetivo de gerenciar tarefas pessoais. Foi criado com foco em consolidar conhecimentos práticos em desenvolvimento backend.

## 🎯 Objetivo
Desenvolver uma aplicação robusta e escalável que permita aos usuários gerenciar suas tarefas diárias, com funcionalidades modernas e boas práticas de desenvolvimento.

## 🚀 Funcionalidades
- Cadastro de usuários
- Autenticação via JWT
- Criação, listagem, edição e exclusão de tarefas
- Marcar tarefas como concluídas
- Categorizar tarefas (importantes, casuais, etc.)
- Compartilhar tarefas com outros usuários
- Alertas de tarefas do dia e/ou importantes

## 🧰 Tecnologias Utilizadas
- **Go**: linguagem principal
- **Gin**: framework web
- **PostgreSQL ou MySQL**: banco de dados relacional
- **Redis**: cache para tarefas por usuário
- **Docker**: ambiente de desenvolvimento
- **Testes unitários**: com `testing` e `testify`
- **Logs**: com `logrus` ou `zap`

## 📦 Estrutura do Projeto
Organizado em módulos como `auth`, `user`, `task`, `middleware`, `cache`, `logger`, seguindo boas práticas de arquitetura em Go.

## 📌 Objetivo Educacional
Este projeto visa aprimorar habilidades técnicas e práticas essenciais como desenvolvedor, incluindo:
- Boas práticas de código
- Testes unitários
- Uso de containers
- Integração com banco de dados e cache
- Versionamento com Git e GitHub


---


# Task Management API (To-Do List) in Go

This project is a RESTful API developed in Go to manage personal tasks. It was created as focusing on consolidating practical backend development skills.

## 🎯 Purpose
Build a robust and scalable application that allows users to manage their daily tasks, with modern features and development best practices.

## 🚀 Features
- User registration
- JWT authentication
- Create, list, edit, and delete tasks
- Mark tasks as completed
- Categorize tasks (important, casual, etc.)
- Share tasks with other users
- Alerts for recent tasks

## 🧰 Technologies Used
- **Go**: main programming language
- **Gin**: web framework
- **PostgreSQL or MySQL**: relational database
- **Redis**: cache for user tasks
- **Docker**: development environment
- **Unit testing**: using `testing` and `testify`
- **Logging**: with `logrus` or `zap`

## 📦 Project Structure
Organized into modules such as `auth`, `user`, `task`, `middleware`, `cache`, `logger`, following Go architecture best practices.

## 📌 Educational Goal
This project aims to improve essential technical and practical skills as a developer, including:
- Clean code practices
- Unit tests
- Container usage
- Integration with database and cache
- Version control with Git and GitHub
