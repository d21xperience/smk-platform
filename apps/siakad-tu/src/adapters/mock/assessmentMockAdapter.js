// apps/siakad-tu/src/adapters/mock/assessmentMockAdapter.js

import { mockStorage } from './mockStorage.js'
import { idGenerator } from '../utils/idGenerator.js'

const _simulateLatency = (ms = 50) => new Promise((resolve) => setTimeout(resolve, ms))
const _error = (code, message, details = null) => ({
  success: false,
  data: null,
  error: { code, message, details },
})
const _success = (data) => ({ success: true, data, error: null })

/**
 * Helper: Konversi nilai angka ke predikat huruf
 * Skala: A (85-100), B (70-84), C (55-69), D (40-54), E (0-39)
 */
const _calculatePredicate = (score) => {
  if (score >= 85) return 'A'
  if (score >= 70) return 'B'
  if (score >= 55) return 'C'
  if (score >= 40) return 'D'
  return 'E'
}

/**
 * Helper: Hitung rata-rata tertimbang
 * Formula: (dailyAvg * 0.4) + (midterm * 0.3) + (final * 0.3)
 * Weight bisa disesuaikan per sekolah
 */
const _calculateWeightedAverage = (dailyAvg, midterm, final) => {
  const DAILY_WEIGHT = 0.4
  const MIDTERM_WEIGHT = 0.3
  const FINAL_WEIGHT = 0.3

  return (
    Math.round(dailyAvg * DAILY_WEIGHT + midterm * MIDTERM_WEIGHT + final * FINAL_WEIGHT * 100) /
    100
  )
}

/**
 * Helper: Hitung grade summary dari daftar assessment
 */
const _computeGradeSummary = (assessments, studentId, subjectId, periodId) => {
  const dailyScores = assessments.filter((a) => a.type === 'DAILY')
  const midtermScores = assessments.filter((a) => a.type === 'MIDTERM')
  const finalScores = assessments.filter((a) => a.type === 'FINAL')

  const dailyAverage =
    dailyScores.length > 0
      ? dailyScores.reduce((sum, a) => sum + (a.score / a.maxScore) * 100, 0) / dailyScores.length
      : 0

  const midtermScore =
    midtermScores.length > 0
      ? midtermScores.reduce((sum, a) => sum + (a.score / a.maxScore) * 100, 0) /
        midtermScores.length
      : 0

  const finalScore =
    finalScores.length > 0
      ? finalScores.reduce((sum, a) => sum + (a.score / a.maxScore) * 100, 0) / finalScores.length
      : 0

  const weightedAverage = _calculateWeightedAverage(
    Math.round(dailyAverage * 100) / 100,
    Math.round(midtermScore * 100) / 100,
    Math.round(finalScore * 100) / 100,
  )

  return {
    studentId,
    subjectId,
    periodId,
    dailyAverage: Math.round(dailyAverage * 100) / 100,
    midtermScore: Math.round(midtermScore * 100) / 100,
    finalScore: Math.round(finalScore * 100) / 100,
    weightedAverage,
    predicate: _calculatePredicate(weightedAverage),
    dailyCount: dailyScores.length,
    isComplete: dailyScores.length > 0 && midtermScores.length > 0 && finalScores.length > 0,
  }
}

export const assessmentMockAdapter = {
  async createAssessment(command, context) {
    await _simulateLatency(100)

    // Validasi: Score tidak boleh melebihi maxScore
    if (command.score > (command.maxScore || 100)) {
      return _error(
        'INVALID_SCORE',
        `Nilai (${command.score}) tidak boleh melebihi skor maksimal (${command.maxScore || 100}).`,
      )
    }

    // Validasi: Score tidak boleh negatif
    if (command.score < 0) {
      return _error('INVALID_SCORE', 'Nilai tidak boleh negatif.')
    }

    // Validasi: Midterm dan Final hanya boleh 1 per siswa per mapel per periode
    const assessments = mockStorage.read(context.schoolId, 'assessments')

    if (command.type === 'MIDTERM' || command.type === 'FINAL') {
      const existing = assessments.find(
        (a) =>
          a.studentId === command.studentId &&
          a.subjectId === command.subjectId &&
          a.periodId === context.periodId &&
          a.type === command.type,
      )
      if (existing) {
        return _error(
          'DUPLICATE_ASSESSMENT',
          `Nilai ${command.type === 'MIDTERM' ? 'UTS' : 'UAS'} sudah ada untuk siswa ini. Gunakan update, bukan create.`,
        )
      }
    }

    const assessmentData = {
      assessmentId: idGenerator.assessmentId(),
      studentId: command.studentId,
      subjectId: command.subjectId,
      periodId: context.periodId,
      type: command.type,
      score: command.score,
      maxScore: command.maxScore || 100,
      weight: command.weight || 1,
      description: command.description || '',
      date: command.date || new Date().toISOString().split('T')[0],
      schoolId: context.schoolId,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    assessments.push(assessmentData)
    mockStorage.write(context.schoolId, 'assessments', assessments)

    return _success(assessmentData)
  },

  async updateAssessment(assessmentId, command, context) {
    await _simulateLatency(100)

    const assessments = mockStorage.read(context.schoolId, 'assessments')
    const index = assessments.findIndex((a) => a.assessmentId === assessmentId)

    if (index === -1) return _error('NOT_FOUND', 'Data penilaian tidak ditemukan.')

    // Validasi: Score tidak boleh melebihi maxScore
    const maxScore = command.maxScore || assessments[index].maxScore
    if (command.score > maxScore) {
      return _error(
        'INVALID_SCORE',
        `Nilai (${command.score}) tidak boleh melebihi skor maksimal (${maxScore}).`,
      )
    }

    if (command.score !== undefined) assessments[index].score = command.score
    if (command.maxScore !== undefined) assessments[index].maxScore = command.maxScore
    if (command.weight !== undefined) assessments[index].weight = command.weight
    if (command.description !== undefined) assessments[index].description = command.description
    if (command.date !== undefined) assessments[index].date = command.date
    assessments[index].updatedAt = new Date().toISOString()

    mockStorage.write(context.schoolId, 'assessments', assessments)

    return _success(assessments[index])
  },

  async deleteAssessment(assessmentId, context) {
    await _simulateLatency(100)

    const assessments = mockStorage.read(context.schoolId, 'assessments')
    const index = assessments.findIndex((a) => a.assessmentId === assessmentId)

    if (index === -1) return _error('NOT_FOUND', 'Data penilaian tidak ditemukan.')

    assessments.splice(index, 1)
    mockStorage.write(context.schoolId, 'assessments', assessments)

    return _success({ deleted: true })
  },

  async getAssessmentById(assessmentId, context) {
    await _simulateLatency(30)

    const assessments = mockStorage.read(context.schoolId, 'assessments')
    const data = assessments.find((a) => a.assessmentId === assessmentId)

    if (!data) return _error('NOT_FOUND', 'Data penilaian tidak ditemukan.')
    return _success(data)
  },

  async getAssessments(context, filters = {}) {
    await _simulateLatency(100)

    let assessments = mockStorage.read(context.schoolId, 'assessments')

    // Filter by studentId
    if (filters.studentId) {
      assessments = assessments.filter((a) => a.studentId === filters.studentId)
    }

    // Filter by subjectId
    if (filters.subjectId) {
      assessments = assessments.filter((a) => a.subjectId === filters.subjectId)
    }

    // Filter by type
    if (filters.type) {
      assessments = assessments.filter((a) => a.type === filters.type)
    }

    // Filter by periodId
    if (filters.periodId) {
      assessments = assessments.filter((a) => a.periodId === filters.periodId)
    }

    // Sort by date descending
    assessments.sort((a, b) => new Date(b.date) - new Date(a.date))

    const page = filters.page || 1
    const limit = filters.limit || 20
    const total = assessments.length
    const items = assessments.slice((page - 1) * limit, page * limit)

    return _success({ items, total, page, limit, totalPages: Math.ceil(total / limit) })
  },

  async getGradeSummary(studentId, subjectId, context) {
    await _simulateLatency(50)

    const assessments = mockStorage.read(context.schoolId, 'assessments')
    const studentSubjectAssessments = assessments.filter(
      (a) =>
        a.studentId === studentId && a.subjectId === subjectId && a.periodId === context.periodId,
    )

    const summary = _computeGradeSummary(
      studentSubjectAssessments,
      studentId,
      subjectId,
      context.periodId,
    )

    return _success(summary)
  },

  async getStudentReportCard(studentId, context) {
    await _simulateLatency(100)

    const assessments = mockStorage.read(context.schoolId, 'assessments')
    const studentAssessments = assessments.filter(
      (a) => a.studentId === studentId && a.periodId === context.periodId,
    )

    // Group by subjectId
    const subjectIds = [...new Set(studentAssessments.map((a) => a.subjectId))]

    const subjects = subjectIds.map((subjectId) => {
      const subjectAssessments = studentAssessments.filter((a) => a.subjectId === subjectId)
      return _computeGradeSummary(subjectAssessments, studentId, subjectId, context.periodId)
    })

    // Hitung overall average
    const completedSubjects = subjects.filter((s) => s.isComplete)
    const overallAverage =
      completedSubjects.length > 0
        ? Math.round(
            (completedSubjects.reduce((sum, s) => sum + s.weightedAverage, 0) /
              completedSubjects.length) *
              100,
          ) / 100
        : 0

    return _success({
      studentId,
      periodId: context.periodId,
      subjects,
      overallAverage,
      overallPredicate: _calculatePredicate(overallAverage),
    })
  },

  async getClassGradeSummary(classId, subjectId, context) {
    await _simulateLatency(100)

    const assessments = mockStorage.read(context.schoolId, 'assessments')
    const students = mockStorage.read(context.schoolId, 'students')

    // Dapatkan semua siswa di kelas ini (dari enrollment aktif)
    const classStudentIds = students
      .filter((s) =>
        s.enrollments?.some(
          (e) => e.classId === classId && e.periodId === context.periodId && e.status === 'active',
        ),
      )
      .map((s) => s.studentId)

    // Filter assessments untuk kelas ini
    const classAssessments = assessments.filter(
      (a) =>
        classStudentIds.includes(a.studentId) &&
        a.subjectId === subjectId &&
        a.periodId === context.periodId,
    )

    // Hitung summary per siswa
    const studentSummaries = classStudentIds.map((studentId) => {
      const studentAssessments = classAssessments.filter((a) => a.studentId === studentId)
      return _computeGradeSummary(studentAssessments, studentId, subjectId, context.periodId)
    })

    // Hitung statistik kelas
    const weightedAverages = studentSummaries.map((s) => s.weightedAverage).filter((a) => a > 0)
    const classAverage =
      weightedAverages.length > 0
        ? Math.round(
            (weightedAverages.reduce((sum, a) => sum + a, 0) / weightedAverages.length) * 100,
          ) / 100
        : 0

    return _success({
      classId,
      subjectId,
      periodId: context.periodId,
      students: studentSummaries,
      classAverage,
      highestScore: weightedAverages.length > 0 ? Math.max(...weightedAverages) : 0,
      lowestScore: weightedAverages.length > 0 ? Math.min(...weightedAverages) : 0,
    })
  },
}
