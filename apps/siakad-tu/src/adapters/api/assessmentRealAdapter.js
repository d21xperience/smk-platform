// apps/siakad-tu/src/adapters/api/assessmentRealAdapter.js

import { api } from '@/boot/axios.js'

/**
 * Helper: Map Backend Assessment DTO ke Frontend Contract
 */
const mapBackendToAssessment = (backendDto) => {
  if (!backendDto) return null

  return {
    assessmentId: backendDto.assessment_id || backendDto.id,
    studentId: backendDto.student_id,
    subjectId: backendDto.subject_id,
    periodId: backendDto.period_id,
    type: backendDto.type, // 'DAILY', 'MIDTERM', 'FINAL'
    score: backendDto.score,
    maxScore: backendDto.max_score || 100,
    weight: backendDto.weight || 1,
    description: backendDto.description,
    date: backendDto.date,
    createdAt: backendDto.created_at || backendDto.createdAt,
    updatedAt: backendDto.updated_at || backendDto.updatedAt,
  }
}

/**
 * Helper: Map Backend Grade Summary DTO ke Frontend Contract
 */
const mapBackendToGradeSummary = (backendDto) => {
  if (!backendDto) return null

  return {
    studentId: backendDto.student_id,
    subjectId: backendDto.subject_id,
    periodId: backendDto.period_id,
    dailyAverage: backendDto.daily_average || 0,
    midtermScore: backendDto.midterm_score || 0,
    finalScore: backendDto.final_score || 0,
    weightedAverage: backendDto.weighted_average || 0,
    predicate: backendDto.predicate || '-', // 'A', 'B', 'C', 'D', 'E'
    dailyCount: backendDto.daily_count || 0,
    isComplete: backendDto.is_complete || false, // Apakah sudah ada semua komponen nilai
  }
}

/**
 * Helper: Konversi ApplicationError ke format Service yang konsisten
 */
const handleApiError = (error) => ({
  success: false,
  data: null,
  error: {
    code: error.code || 'API_ERROR',
    message: error.message || 'Terjadi kesalahan pada server.',
    details: error.details || null,
  },
})

const handleSuccess = (data) => ({
  success: true,
  data,
  error: null,
})

export const assessmentRealAdapter = {
  /**
   * Input nilai baru (Daily/Midterm/Final)
   */
  async createAssessment(command, context) {
    try {
      const payload = {
        studentId: command.studentId,
        subjectId: command.subjectId,
        type: command.type,
        score: command.score,
        maxScore: command.maxScore || 100,
        weight: command.weight || 1,
        description: command.description,
        date: command.date,
        periodId: context.periodId,
      }
      const response = await api.post('/assessments', payload)
      return handleSuccess(mapBackendToAssessment(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Update nilai yang sudah ada
   */
  async updateAssessment(assessmentId, command) {
    try {
      const payload = {
        score: command.score,
        maxScore: command.maxScore,
        weight: command.weight,
        description: command.description,
        date: command.date,
      }
      const response = await api.put(`/assessments/${assessmentId}`, payload)
      return handleSuccess(mapBackendToAssessment(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Hapus nilai
   */
  async deleteAssessment(assessmentId) {
    try {
      await api.delete(`/assessments/${assessmentId}`)
      return handleSuccess({ deleted: true })
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil detail nilai by ID
   */
  async getAssessmentById(assessmentId) {
    try {
      const response = await api.get(`/assessments/${assessmentId}`)
      return handleSuccess(mapBackendToAssessment(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil daftar nilai dengan filter
   */
  async getAssessments(context, filters = {}) {
    try {
      const params = { schoolId: context.schoolId, periodId: context.periodId, ...filters }
      const response = await api.get('/assessments', { params })

      const mappedData = response.data.items
        ? { ...response.data, items: response.data.items.map(mapBackendToAssessment) }
        : response.data

      return handleSuccess(mappedData)
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil ringkasan nilai siswa per mata pelajaran
   */
  async getGradeSummary(studentId, subjectId, context) {
    try {
      const response = await api.get(`/students/${studentId}/grade-summary`, {
        params: { subjectId, periodId: context.periodId },
      })
      return handleSuccess(mapBackendToGradeSummary(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil rapor lengkap siswa (semua mata pelajaran)
   */
  async getStudentReportCard(studentId, context) {
    try {
      const response = await api.get(`/students/${studentId}/report-card`, {
        params: { periodId: context.periodId },
      })
      return handleSuccess({
        studentId: response.data.student_id,
        periodId: response.data.period_id,
        subjects: response.data.subjects.map(mapBackendToGradeSummary),
        overallAverage: response.data.overall_average || 0,
        overallPredicate: response.data.overall_predicate || '-',
      })
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil rekap nilai per kelas per mata pelajaran
   */
  async getClassGradeSummary(classId, subjectId, context) {
    try {
      const response = await api.get(`/classes/${classId}/grade-summary`, {
        params: { subjectId, periodId: context.periodId },
      })
      return handleSuccess({
        classId,
        subjectId,
        periodId: context.periodId,
        students: response.data.students.map(mapBackendToGradeSummary),
        classAverage: response.data.class_average || 0,
        highestScore: response.data.highest_score || 0,
        lowestScore: response.data.lowest_score || 0,
      })
    } catch (error) {
      return handleApiError(error)
    }
  },
}
