import axios from 'axios'

const http = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:3000',
  withCredentials: false,
})

export { http }
