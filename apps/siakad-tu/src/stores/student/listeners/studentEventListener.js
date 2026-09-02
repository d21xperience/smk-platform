// apps/siakad-tu/src/stores/student/listeners/studentEventListener.js

import { eventDispatcher } from '../../../events/eventDispatcher.js'
import { useStudentListStore } from '../studentListStore.js'
import { useStudentDetailStore } from '../studentDetailStore.js'
import { useStudentFormStore } from '../studentFormStore.js'

/**
 * Student Event Listener — Mendengarkan Domain Events dan meresponsnya
 * dengan mengupdate state store.
 *
 * Ini adalah implementasi Projection Pattern di frontend:
 * Domain Event → Update Read Model (Store)
 *
 * Dipanggil sekali saat app init (di boot file).
 */
class StudentEventListener {
  constructor() {
    this._initialized = false
  }

  /**
   * Inisialisasi listener. Panggil sekali di boot file.
   */
  init() {
    if (this._initialized) return

    // === Context Change Listener ===
    // Saat user ganti sekolah/periode, reset semua store
    eventDispatcher.on('OperationalContextChanged', () => {
      this._handleContextChanged()
    })

    eventDispatcher.on('OperationalContextCleared', () => {
      this._handleContextCleared()
    })

    // === Student Domain Events ===
    // Update cache saat ada perubahan data

    eventDispatcher.on('StudentRegistered', (event) => {
      this._handleStudentRegistered(event)
    })

    eventDispatcher.on('StudentEnrolled', (event) => {
      this._handleStudentEnrolled(event)
    })

    eventDispatcher.on('StudentProfileUpdated', (event) => {
      this._handleStudentProfileUpdated(event)
    })

    eventDispatcher.on('StudentGraduated', (event) => {
      this._handleStudentGraduated(event)
    })

    eventDispatcher.on('StudentTransferred', (event) => {
      this._handleStudentTransferred(event)
    })

    this._initialized = true
    console.log('[StudentEventListener] Initialized.')
  }

  /**
   * Handle context changed — reset semua store
   * @private
   */
  _handleContextChanged() {
    console.log('[StudentEventListener] Context changed, resetting stores...')

    // Reset stores
    try {
      const listStore = useStudentListStore()
      const detailStore = useStudentDetailStore()
      const formStore = useStudentFormStore()

      listStore.reset()
      detailStore.reset()
      formStore.reset()
    } catch (error) {
      // Store mungkin belum di-inisialisasi (misal sebelum login)
      console.warn('[StudentEventListener] Could not reset stores:', error.message)
    }
  }

  /**
   * Handle context cleared (logout)
   * @private
   */
  _handleContextCleared() {
    this._handleContextChanged()
  }

  /**
   * Handle StudentRegistered — invalidate list cache
   * @private
   */
  _handleStudentRegistered(event) {
    console.log('[StudentEventListener] StudentRegistered:', event.payload)
    try {
      const listStore = useStudentListStore()
      // Force refetch list di next access
      listStore.initialized = false
    } catch (error) {
      console.warn('[StudentEventListener] Error handling StudentRegistered:', error.message)
    }
  }

  /**
   * Handle StudentEnrolled — invalidate detail cache
   * @private
   */
  _handleStudentEnrolled(event) {
    console.log('[StudentEventListener] StudentEnrolled:', event.payload)
    try {
      const detailStore = useStudentDetailStore()
      detailStore.invalidate(event.aggregateId)
    } catch (error) {
      console.warn('[StudentEventListener] Error handling StudentEnrolled:', error.message)
    }
  }

  /**
   * Handle StudentProfileUpdated — invalidate detail cache
   * @private
   */
  _handleStudentProfileUpdated(event) {
    console.log('[StudentEventListener] StudentProfileUpdated:', event.payload)
    try {
      const detailStore = useStudentDetailStore()
      detailStore.invalidate(event.aggregateId)
    } catch (error) {
      console.warn('[StudentEventListener] Error handling StudentProfileUpdated:', error.message)
    }
  }

  /**
   * Handle StudentGraduated — remove from active list
   * @private
   */
  _handleStudentGraduated(event) {
    console.log('[StudentEventListener] StudentGraduated:', event.payload)
    try {
      const listStore = useStudentListStore()
      const detailStore = useStudentDetailStore()

      listStore.removeItem(event.aggregateId)
      detailStore.invalidate(event.aggregateId)
    } catch (error) {
      console.warn('[StudentEventListener] Error handling StudentGraduated:', error.message)
    }
  }

  /**
   * Handle StudentTransferred — remove from list
   * @private
   */
  _handleStudentTransferred(event) {
    console.log('[StudentEventListener] StudentTransferred:', event.payload)
    try {
      const listStore = useStudentListStore()
      const detailStore = useStudentDetailStore()

      listStore.removeItem(event.aggregateId)
      detailStore.invalidate(event.aggregateId)
    } catch (error) {
      console.warn('[StudentEventListener] Error handling StudentTransferred:', error.message)
    }
  }
}

// Singleton instance
export const studentEventListener = new StudentEventListener()
