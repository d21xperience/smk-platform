import { defineStore } from 'pinia'
import { ref } from 'vue'
import { letterRealAdapter } from '@/adapters/api/letterRealAdapter.js'

export const useLetterStore = defineStore('letter', () => {
  const letters = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  async function fetchLetters() {
    isLoading.value = true
    try {
      const res = await letterRealAdapter.getLetters()
      if (res.success) letters.value = res.data
      else error.value = res.error
    } catch (err) {
      error.value = err.message
    } finally {
      isLoading.value = false
    }
  }

  async function createLetter(payload) {
    isLoading.value = true
    try {
      const res = await letterRealAdapter.createLetter(payload)
      if (res.success) {
        await fetchLetters()
        return true
      }
      error.value = res.error
      return false
    } catch (err) {
      error.value = err.message
      return false
    } finally {
      isLoading.value = false
    }
  }

  async function completeLetter(id) {
    isLoading.value = true
    try {
      const res = await letterRealAdapter.completeLetter(id)
      if (res.success) {
        await fetchLetters()
        return true
      }
      return false
    } finally {
      isLoading.value = false
    }
  }

  return { letters, isLoading, error, fetchLetters, createLetter, completeLetter }
})
