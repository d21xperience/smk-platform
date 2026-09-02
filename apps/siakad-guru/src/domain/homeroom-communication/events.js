// FILE: src/domain/homeroom-communication/events.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Communication Domain Events
 * Event types dan factory functions untuk publish melalui EventDispatcher.
 * Events ini dipublish oleh Service setelah mutasi berhasil.
 */

export const HOMEROOM_COMMUNICATION_EVENTS = {
  COMMUNICATION_CREATED: 'homeroom.communication.created',
  COMMUNICATION_STATUS_UPDATED: 'homeroom.communication.status.updated',
}

/**
 * Factory untuk event Communication Created.
 * @param {Object} communication - Log komunikasi yang baru dibuat.
 * @param {string} userId - User ID yang membuat.
 * @returns {Object} Event payload.
 */
export function createCommunicationCreatedEvent(communication, userId) {
  return {
    type: HOMEROOM_COMMUNICATION_EVENTS.COMMUNICATION_CREATED,
    payload: {
      communicationId: communication.communicationId,
      studentId: communication.studentId,
      studentName: communication.studentName,
      category: communication.category,
      direction: communication.direction,
      title: communication.title,
      status: communication.status,
      communicationDate: communication.communicationDate,
      createdBy: userId,
      timestamp: new Date().toISOString(),
    },
  }
}

/**
 * Factory untuk event Communication Status Updated.
 * @param {Object} communication - Log komunikasi yang diupdate.
 * @param {string} previousStatus - Status sebelumnya.
 * @param {string} userId - User ID yang mengupdate.
 * @returns {Object} Event payload.
 */
export function createCommunicationStatusUpdatedEvent(communication, previousStatus, userId) {
  return {
    type: HOMEROOM_COMMUNICATION_EVENTS.COMMUNICATION_STATUS_UPDATED,
    payload: {
      communicationId: communication.communicationId,
      studentId: communication.studentId,
      studentName: communication.studentName,
      title: communication.title,
      previousStatus: previousStatus,
      newStatus: communication.status,
      updatedBy: userId,
      timestamp: new Date().toISOString(),
    },
  }
}
