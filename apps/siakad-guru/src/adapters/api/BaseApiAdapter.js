import { api } from '../../boot/axios.js'

export class BaseApiAdapter {
  constructor(basePath) {
    this.basePath = basePath
  }

  async get(endpoint, params = {}) {
    try {
      const response = await api.get(`${this.basePath}${endpoint}`, { params })
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  }

  async post(endpoint, data = {}) {
    try {
      const response = await api.post(`${this.basePath}${endpoint}`, data)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  }

  async put(endpoint, data = {}) {
    try {
      const response = await api.put(`${this.basePath}${endpoint}`, data)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  }

  async delete(endpoint) {
    try {
      const response = await api.delete(`${this.basePath}${endpoint}`)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  }

  handleError(error) {
    if (error.response) {
      const { status, data } = error.response
      const message = data?.message || data?.error || 'An error occurred'

      switch (status) {
        case 400:
          return new Error(`Bad Request: ${message}`)
        case 401:
          return new Error(`Unauthorized: ${message}`)
        case 403:
          return new Error(`Forbidden: ${message}`)
        case 404:
          return new Error(`Not Found: ${message}`)
        case 422:
          return new Error(`Validation Error: ${message}`)
        case 500:
          return new Error(`Server Error: ${message}`)
        default:
          return new Error(`Error ${status}: ${message}`)
      }
    }

    if (error.request) {
      return new Error('Network error: Unable to reach server')
    }

    return error
  }
}
