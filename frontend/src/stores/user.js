import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: null
  }),
  
  actions: {
    async login(username, password) {
      try {
        const response = await axios.post('http://localhost:8080/api/v1/users/login', {
          username,
          password
        })
        const token = response.data.data
        this.token = token
        localStorage.setItem('token', token)
        return true
      } catch (error) {
        console.error('Login failed:', error)
        return false
      }
    },

    async register(username, password) {
      try {
        await axios.post('http://localhost:8080/api/v1/users/register', {
          username,
          password
        })
        return true
      } catch (error) {
        console.error('Registration failed:', error)
        return false
      }
    },

    async getUserInfo() {
      try {
        const response = await axios.get('http://localhost:8080/api/v1/users/info', {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        this.userInfo = response.data.data
        return true
      } catch (error) {
        console.error('Failed to get user info:', error)
        return false
      }
    },

    logout() {
      this.token = ''
      this.userInfo = null
      localStorage.removeItem('token')
    }
  }
}) 