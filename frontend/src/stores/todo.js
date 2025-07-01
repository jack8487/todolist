import { defineStore } from 'pinia'
import axios from 'axios'
import { useUserStore } from './user'

export const useTodoStore = defineStore('todo', {
  state: () => ({
    todos: []
  }),

  actions: {
    async fetchTodos() {
      try {
        const userStore = useUserStore()
        const response = await axios.get('http://localhost:8080/api/v1/tasks', {
          headers: { Authorization: `Bearer ${userStore.token}` }
        })
        this.todos = response.data.data.items || []
        return true
      } catch (error) {
        console.error('Failed to fetch todos:', error)
        return false
      }
    },

    async addTodo(title, description, dueDate) {
      try {
        const userStore = useUserStore()
        const requestData = {
          title,
          description,
          due_date: dueDate || ''
        }
        console.log('发送的任务数据:', requestData)
        
        const response = await axios.post('http://localhost:8080/api/v1/tasks', requestData, {
          headers: { Authorization: `Bearer ${userStore.token}` }
        })
        console.log('服务器响应:', response.data)
        
        if (response.data.code === 200) {
          const newTask = response.data.data
          this.todos.unshift(newTask)
          return true
        }
        return false
      } catch (error) {
        console.error('添加任务失败:', error.response?.data || error.message)
        return false
      }
    },

    async updateTodo(id, updates) {
      try {
        const userStore = useUserStore()
        const response = await axios.put(`http://localhost:8080/api/v1/tasks/${id}`, updates, {
          headers: { Authorization: `Bearer ${userStore.token}` }
        })
        const index = this.todos.findIndex(todo => todo.id === id)
        if (index !== -1) {
          this.todos[index] = { ...this.todos[index], ...response.data.data }
        }
        return true
      } catch (error) {
        console.error('Failed to update todo:', error)
        return false
      }
    },

    async deleteTodo(id) {
      try {
        const userStore = useUserStore()
        await axios.delete(`http://localhost:8080/api/v1/tasks/${id}`, {
          headers: { Authorization: `Bearer ${userStore.token}` }
        })
        this.todos = this.todos.filter(todo => todo.id !== id)
        return true
      } catch (error) {
        console.error('Failed to delete todo:', error)
        return false
      }
    },

    async toggleTodoStatus(todo) {
      try {
        const userStore = useUserStore()
        const newStatus = todo.status === 'done' ? 'todo' : 'done'
        
        const response = await axios.put(`http://localhost:8080/api/v1/tasks/${todo.id}`, {
          status: newStatus
        }, {
          headers: { Authorization: `Bearer ${userStore.token}` }
        })

        if (response.data.code === 200) {
          // 更新本地任务状态
          const index = this.todos.findIndex(t => t.id === todo.id)
          if (index !== -1) {
            this.todos[index] = response.data.data
          }
          return true
        }
        return false
      } catch (error) {
        console.error('更新任务状态失败:', error.response?.data || error.message)
        return false
      }
    }
  }
}) 