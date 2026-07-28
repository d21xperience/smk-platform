// src/services/AssessmentService.js

export class AssessmentService {
  /**
   * @param {Object} adapter - kontrak adapter (mock/api)
   */
  constructor(adapter) {
    this.adapter = adapter
  }

  /**
   * Ambil atau buat AssessmentSession untuk sesi mengajar tertentu.
   * @param {number} teachingSessionId
   * @param {Object} options - { className, subject, academicYearId, semesterId }
   * @returns {Promise<import('src/models/AssessmentSession').AssessmentSession>}
   */
  async fetchOrCreateSession(teachingSessionId, options) {
    return this.adapter.fetchOrCreateAssessmentSession(teachingSessionId, options)
  }

  /**
   * Simpan AssessmentSession sebagai draft.
   * @param {Object} session - AssessmentSession
   * @returns {Promise<import('src/models/AssessmentSession').AssessmentSession>}
   */
  async saveDraft(session) {
    return this.adapter.saveAssessmentDraft(session)
  }

  /**
   * Submit final assessment.
   * @param {number} sessionId
   * @returns {Promise<import('src/models/AssessmentSession').AssessmentSession>}
   */
  async submit(sessionId) {
    return this.adapter.submitAssessment(sessionId)
  }

  /**
   * Ambil komponen penilaian untuk session.
   * @param {number} sessionId
   * @returns {Promise<Array<import('src/models/AssessmentComponent').AssessmentComponent>>}
   */
  async fetchComponents(sessionId) {
    return this.adapter.getAssessmentComponents(sessionId)
  }

  /**
   * Simpan atau update skor untuk satu komponen.
   * @param {number} componentId
   * @param {Array} scores - array of ScoreEntry
   * @returns {Promise<import('src/models/AssessmentComponent').AssessmentComponent>}
   */
  async saveScores(componentId, scores) {
    return this.adapter.saveComponentScores(componentId, scores)
  }

  /**
   * Ambil hasil akhir (grades) untuk session.
   * @param {number} sessionId
   * @returns {Promise<import('src/models/AssessmentFinalResult').AssessmentFinalResult|null>}
   */
  async fetchFinalResult(sessionId) {
    return this.adapter.getAssessmentFinalResult(sessionId)
  }

  /**
   * Ambil daftar siswa untuk kelas tertentu (jika diperlukan dari luar).
   * @param {string} className
   * @returns {Promise<Array<{ studentId: number|string, studentName: string }>>}
   */
  async fetchStudents(className) {
    return this.adapter.getStudentsByClass(className)
  }
}
