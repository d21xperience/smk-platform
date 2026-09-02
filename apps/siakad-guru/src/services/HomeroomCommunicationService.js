// FILE: src/services/HomeroomCommunicationService.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { HomeroomCommunicationEngine } from '../domain/homeroom-communication/HomeroomCommunicationEngine.js'
import {
  createCommunicationCreatedEvent,
  createCommunicationStatusUpdatedEvent,
} from '../domain/homeroom-communication/events.js'
import {
  mapHomeroomCommunications,
  mapHomeroomCommunication,
  mapHomeroomCommunicationSummary,
} from '../contracts/homeroomCommunicationContract.js'

/**
 * Homeroom Communication Application Service
 * Orchestration use case log komunikasi wali kelas dengan orang tua.
 * Jalur dengan Business Rule + Mutasi:
 * Service memanggil Domain Engine untuk validasi & state transition,
 * Adapter untuk persistence, dan EventDispatcher untuk publish events.
 * Tidak ada operasi delete — log komunikasi tidak dihapus.
 */
export class HomeroomCommunicationService {
  constructor({ homeroomCommunicationAdapter, eventDispatcher }) {
    this.adapter = homeroomCommunicationAdapter
    this.eventDispatcher = eventDispatcher
    this.engine = new HomeroomCommunicationEngine()
  }

  /**
   * Mengambil daftar log komunikasi wali kelas beserta summary.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Object>} { communications, summary }
   */
  async getCommunications(context) {
    if (!context.userId) {
      throw new Error('userId is required to fetch homeroom communications.')
    }
    if (!context.academicYearId) {
      throw new Error('academicYearId is required to fetch homeroom communications.')
    }
    if (!context.semesterId) {
      throw new Error('semesterId is required to fetch homeroom communications.')
    }

    const rawCommunications = await this.adapter.fetchCommunicationsByHomeroom(context)
    const communications = mapHomeroomCommunications(rawCommunications)
    const sortedCommunications = this.engine.sortCommunications(communications)
    const summary = mapHomeroomCommunicationSummary(this.engine.buildSummary(sortedCommunications))

    return { communications: sortedCommunications, summary }
  }

  /**
   * Membuat log komunikasi baru.
   * Flow: Engine.validate → Engine.create → Adapter.save → Dispatch event.
   * @param {Object} context - Operational context
   * @param {Object} commData - Data log komunikasi dari form
   * @returns {Promise<Object>} Log komunikasi yang tersimpan
   */
  async createCommunication(context, commData) {
    if (!context.userId) {
      throw new Error('userId is required to create communication.')
    }

    this.engine.validateCommunicationInput(commData)

    const communication = this.engine.createCommunication(commData, context.userId)

    const savedCommunication = await this.adapter.saveCommunication(communication, context)

    const mappedCommunication = mapHomeroomCommunication(savedCommunication)

    this.eventDispatcher.dispatch(
      createCommunicationCreatedEvent(mappedCommunication, context.userId),
    )

    return mappedCommunication
  }

  /**
   * Mengupdate status tindak lanjut log komunikasi.
   * Flow: Engine.validateTransition → Engine.applyChange → Adapter.update → Dispatch event.
   * @param {Object} context - Operational context
   * @param {Object} communication - Log komunikasi yang akan diupdate
   * @param {string} newStatus - Status baru (FOLLOWED_UP atau CLOSED)
   * @param {string} followUpNote - Catatan tindak lanjut (opsional)
   * @returns {Promise<Object>} Log komunikasi yang terupdate
   */
  async updateStatus(context, communication, newStatus, followUpNote) {
    if (!context.userId) {
      throw new Error('userId is required to update communication status.')
    }

    const previousStatus = communication.status

    this.engine.validateStatusTransition(communication, newStatus)

    const updatedCommunication = this.engine.applyStatusChange(
      communication,
      newStatus,
      followUpNote,
    )

    const savedCommunication = await this.adapter.updateCommunication(updatedCommunication, context)

    const mappedCommunication = mapHomeroomCommunication(savedCommunication)

    this.eventDispatcher.dispatch(
      createCommunicationStatusUpdatedEvent(mappedCommunication, previousStatus, context.userId),
    )

    return mappedCommunication
  }

  /**
   * Bangun summary dari array log komunikasi.
   * Digunakan oleh Store untuk rebuild summary setelah mutasi lokal.
   * @param {Array<Object>} communications - Array log komunikasi.
   * @returns {Object} Summary yang telah dipetakan ke kontrak UI.
   */
  buildSummary(communications) {
    return mapHomeroomCommunicationSummary(this.engine.buildSummary(communications))
  }
}
