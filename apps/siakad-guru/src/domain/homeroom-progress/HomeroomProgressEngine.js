// FILE: src/domain/homeroom-progress/HomeroomProgressEngine.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Progress Engine
 * Pure JavaScript Domain Engine untuk business rules rekap nilai wali kelas.
 * TIDAK BOLEH import: Vue, Pinia, Axios, Quasar, HTTP, Database, Browser API.
 *
 * Business Rules:
 * - Bobot nilai akhir: 40% formatif + 60% sumatif
 * - Status KKM: finalScore >= kkm → PASS, else FAIL
 * - Predikat: A >= 90, B >= 80, C >= 70, D < 70
 * - Tren: UP jika naik >= 5, DOWN jika turun >= 5, STABLE jika selisih < 5, NEW jika tidak ada data sebelumnya
 * - Ranking: Sort by overallAverage DESC
 * - Validasi: Score harus 0-100
 */

const FORMATIVE_WEIGHT = 0.4
const SUMMATIVE_WEIGHT = 0.6
const TREND_THRESHOLD = 5

export const SCORE_STATUS = {
  PASS: 'PASS',
  FAIL: 'FAIL',
}

export const SCORE_TREND = {
  UP: 'UP',
  DOWN: 'DOWN',
  STABLE: 'STABLE',
  NEW: 'NEW',
}

export const SCORE_PREDICATE = {
  A: 'A',
  B: 'B',
  C: 'C',
  D: 'D',
}

export class HomeroomProgressEngine {
  /**
   * Validasi rentang skor (0-100).
   * Jika null/undefined/NaN → return 0.
   * Jika di luar rentang → clamp ke 0 atau 100.
   * @param {number} score - Skor mentah.
   * @returns {number} Skor valid (0-100).
   */
  validateScore(score) {
    if (score === null || score === undefined) return 0
    const num = Number(score)
    if (Number.isNaN(num)) return 0
    return Math.max(0, Math.min(100, num))
  }

  /**
   * Hitung nilai akhir berdasarkan bobot formatif dan sumatif.
   * Formula: (formative * 0.4) + (summative * 0.6)
   * @param {number} formativeScore - Nilai formatif.
   * @param {number} summativeScore - Nilai sumatif.
   * @returns {number} Nilai akhir (dibulatkan 1 desimal).
   */
  calculateFinalScore(formativeScore, summativeScore) {
    const validFormative = this.validateScore(formativeScore)
    const validSummative = this.validateScore(summativeScore)
    const finalScore = (validFormative * FORMATIVE_WEIGHT) + (validSummative * SUMMATIVE_WEIGHT)
    return Math.round(finalScore * 10) / 10
  }

  /**
   * Tentukan status KKM.
   * @param {number} finalScore - Nilai akhir.
   * @param {number} kkm - KKM mata pelajaran.
   * @returns {string} 'PASS' atau 'FAIL'.
   */
  determineStatus(finalScore, kkm) {
    const validKkm = this.validateScore(kkm)
    return finalScore >= validKkm ? SCORE_STATUS.PASS : SCORE_STATUS.FAIL
  }

  /**
   * Tentukan predikat berdasarkan nilai akhir.
   * A >= 90, B >= 80, C >= 70, D < 70
   * @param {number} finalScore - Nilai akhir.
   * @returns {string} Predikat ('A', 'B', 'C', 'D').
   */
  determinePredicate(finalScore) {
    if (finalScore >= 90) return SCORE_PREDICATE.A
    if (finalScore >= 80) return SCORE_PREDICATE.B
    if (finalScore >= 70) return SCORE_PREDICATE.C
    return SCORE_PREDICATE.D
  }

  /**
   * Hitung tren perkembangan berdasarkan perbandingan dengan nilai semester lalu.
   * UP: naik >= 5 poin
   * DOWN: turun >= 5 poin
   * STABLE: selisih < 5 poin
   * NEW: tidak ada data sebelumnya
   * @param {number} currentFinalScore - Nilai akhir semester ini.
   * @param {number|null} previousFinalScore - Nilai akhir semester lalu (null jika tidak ada).
   * @returns {string} Tren ('UP', 'DOWN', 'STABLE', 'NEW').
   */
  calculateTrend(currentFinalScore, previousFinalScore) {
    if (previousFinalScore === null || previousFinalScore === undefined) {
      return SCORE_TREND.NEW
    }

    const diff = currentFinalScore - previousFinalScore
    if (diff >= TREND_THRESHOLD) return SCORE_TREND.UP
    if (diff <= -TREND_THRESHOLD) return SCORE_TREND.DOWN
    return SCORE_TREND.STABLE
  }

  /**
   * Hitung rata-rata keseluruhan dari array nilai.
   * @param {Array<number>} scores - Array nilai.
   * @returns {number} Rata-rata (dibulatkan 1 desimal).
   */
  calculateAverage(scores) {
    if (!Array.isArray(scores) || scores.length === 0) return 0
    const total = scores.reduce((sum, score) => sum + score, 0)
    return Math.round((total / scores.length) * 10) / 10
  }

  /**
   * Hitung ranking kelas berdasarkan overall average.
   * Siswa dengan rata-rata tertinggi mendapat rank 1.
   * Jika ada nilai sama, mendapat rank yang sama.
   * @param {Array<Object>} students - Array siswa dengan property overallAverage.
   * @returns {Array<Object>} Array siswa dengan property classRank ditambahkan.
   */
  calculateClassRank(students) {
    if (!Array.isArray(students) || students.length === 0) return []

    const sorted = [...students].sort((a, b) => b.overallAverage - a.overallAverage)

    let currentRank = 1
    for (let i = 0; i < sorted.length; i++) {
      if (i > 0 && sorted[i].overallAverage < sorted[i - 1].overallAverage) {
        currentRank = i + 1
      }
      sorted[i].classRank = currentRank
    }

    return sorted
  }

  /**
   * Proses data mentah satu siswa menjadi rekap lengkap.
   * @param {Object} rawStudent - Data mentah siswa dari adapter.
   * @param {Array<Object>} subjects - Daftar mata pelajaran dengan KKM.
   * @returns {Object} Rekap siswa yang telah diproses.
   */
  processStudent(rawStudent, subjects) {
    const processedSubjects = []
    const finalScores = []
    let passCount = 0
    let failCount = 0

    for (const subject of subjects) {
      const subjectId = subject.subjectId
      const scoreData = rawStudent.scores[subjectId] || {}

      const formativeScore = this.validateScore(scoreData.formative)
      const summativeScore = this.validateScore(scoreData.summative)
      const previousFinalScore = scoreData.previousFinalScore !== undefined
        ? scoreData.previousFinalScore
        : null

      const finalScore = this.calculateFinalScore(formativeScore, summativeScore)
      const kkm = this.validateScore(subject.kkm)
      const status = this.determineStatus(finalScore, kkm)
      const trend = this.calculateTrend(finalScore, previousFinalScore)
      const predicate = this.determinePredicate(finalScore)

      finalScores.push(finalScore)

      if (status === SCORE_STATUS.PASS) {
        passCount++
      } else {
        failCount++
      }

      processedSubjects.push({
        subjectId: subjectId,
        subjectName: subject.subjectName,
        formativeScore: formativeScore,
        summativeScore: summativeScore,
        finalScore: finalScore,
        kkm: kkm,
        status: status,
        trend: trend,
        predicate: predicate,
      })
    }

    const overallAverage = this.calculateAverage(finalScores)

    return {
      studentId: rawStudent.studentId,
      studentName: rawStudent.studentName,
      seatNumber: rawStudent.seatNumber || 0,
      subjects: processedSubjects,
      overallAverage: overallAverage,
      classRank: 0,
      totalSubjects: subjects.length,
      passCount: passCount,
      failCount: failCount,
    }
  }

  /**
   * Proses seluruh data mentah kelas menjadi rekap lengkap dengan ranking dan summary.
   * @param {Object} rawData - Data mentah dari adapter.
   * @param {Array<Object>} rawData.students - Array data siswa mentah.
   * @param {Array<Object>} rawData.subjects - Array mata pelajaran dengan KKM.
   * @param {string} rawData.classId - ID kelas.
   * @param {string} rawData.className - Nama kelas.
   * @returns {Object} Rekap lengkap kelas.
   */
  processClassProgress(rawData) {
    const { students, subjects, classId, className } = rawData

    if (!Array.isArray(students) || !Array.isArray(subjects)) {
      return {
        classId: classId || null,
        className: className || null,
        summary: {
          totalStudents: 0,
          classAverage: 0,
          highestAverage: 0,
          lowestAverage: 0,
          passRate: 0,
        },
        students: [],
      }
    }

    const processedStudents = students.map(student => this.processStudent(student, subjects))

    const rankedStudents = this.calculateClassRank(processedStudents)

    const allAverages = rankedStudents.map(s => s.overallAverage)
    const classAverage = this.calculateAverage(allAverages)
    const highestAverage = allAverages.length > 0 ? Math.max(...allAverages) : 0
    const lowestAverage = allAverages.length > 0 ? Math.min(...allAverages) : 0

    const totalPass = rankedStudents.reduce((sum, s) => sum + s.passCount, 0)
    const totalSubjectsAll = rankedStudents.reduce((sum, s) => sum + s.totalSubjects, 0)
    const passRate = totalSubjectsAll > 0
      ? Math.round((totalPass / totalSubjectsAll) * 100)
      : 0

    return {
      classId: classId || null,
      className: className || null,
      summary: {
        totalStudents: rankedStudents.length,
        classAverage: classAverage,
        highestAverage: highestAverage,
        lowestAverage: lowestAverage,
        passRate: passRate,
      },
      students: rankedStudents,
    }
  }
}
