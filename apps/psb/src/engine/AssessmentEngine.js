// src/engine/AssessmentEngine.js

export const ComponentType = {
  PH: 'PH',
  PTS: 'PTS',
  PAS: 'PAS',
  TUGAS: 'TUGAS',
  PRAKTIK: 'PRAKTIK',
  PROYEK: 'PROYEK',
  PORTOFOLIO: 'PORTOFOLIO',
  SIKAP: 'SIKAP',
}

// Predikat default (dapat dikustomisasi)
export const DefaultPredicates = [
  { letter: 'A', minScore: 90, maxScore: 100, description: 'Sangat Baik' },
  { letter: 'B', minScore: 75, maxScore: 89, description: 'Baik' },
  { letter: 'C', minScore: 60, maxScore: 74, description: 'Cukup' },
  { letter: 'D', minScore: 0, maxScore: 59, description: 'Kurang' },
]

export class AssessmentEngine {
  /**
   * Hitung nilai akhir seorang siswa dari komponen & skor.
   * Rumus: Σ( (score / maxScore) * weight ) / totalWeight * 100
   * @param {Array<{ weight: number, maxScore: number }>} components
   * @param {Array<{ score: number }>} scores - paralel dengan components
   * @returns {number} nilai akhir (0-100)
   */
  static calculateFinalScore(components, scores) {
    if (components.length === 0) return 0
    let totalWeight = 0
    let weightedSum = 0
    for (let i = 0; i < components.length; i++) {
      const comp = components[i]
      const scoreEntry = scores[i]
      if (!scoreEntry || scoreEntry.score === undefined || scoreEntry.score === null) continue
      const weight = comp.weight || 0
      const maxScore = comp.maxScore || 1
      const normalized = scoreEntry.score / maxScore
      weightedSum += normalized * weight
      totalWeight += weight
    }
    if (totalWeight === 0) return 0
    return Math.round((weightedSum / totalWeight) * 100 * 100) / 100 // 2 desimal
  }

  /**
   * Konversi nilai akhir ke predikat.
   * @param {number} finalScore
   * @param {Array<{ letter: string, minScore: number, maxScore: number }>} [predicates=DefaultPredicates]
   * @returns {{ letter: string, description: string }}
   */
  static getPredicate(finalScore, predicates = DefaultPredicates) {
    for (const p of predicates) {
      if (finalScore >= p.minScore && finalScore <= p.maxScore) {
        return { letter: p.letter, description: p.description }
      }
    }
    return { letter: '-', description: 'Tidak Ada Predikat' }
  }

  /**
   * Validasi apakah session bisa disubmit.
   * Syarat:
   * - Semua komponen wajib (isMandatory) harus memiliki skor untuk semua siswa.
   * - Tidak boleh ada skor kosong (null/undefined) pada komponen wajib.
   * @param {Object} session - AssessmentSession
   * @returns {{ valid: boolean, errors: string[] }}
   */
  static validateForSubmit(session) {
    const errors = []
    if (!session.components || session.components.length === 0) {
      errors.push('Tidak ada komponen penilaian.')
      return { valid: false, errors }
    }
    const studentIds = new Set()
    session.components.forEach((comp) => {
      comp.scores.forEach((s) => studentIds.add(s.studentId))
    })

    for (const comp of session.components) {
      if (!comp.isMandatory) continue
      // Setiap siswa wajib punya skor di komponen ini
      for (const sid of studentIds) {
        const scoreEntry = comp.scores.find((s) => s.studentId === sid)
        if (!scoreEntry || scoreEntry.score === undefined || scoreEntry.score === null) {
          errors.push(`Komponen '${comp.name}' wajib diisi untuk semua siswa.`)
          break // satu error per komponen cukup
        }
      }
    }
    if (session.status === 'submitted') {
      errors.push('Assessment sudah disubmit.')
    }
    return { valid: errors.length === 0, errors }
  }

  /**
   * Menghasilkan daftar Grade untuk semua siswa berdasarkan session.
   * @param {Object} session - AssessmentSession
   * @param {Array<{ letter: string, minScore: number, maxScore: number }>} [predicates]
   * @returns {Array<{ studentId, studentName, finalScore, predicate, componentDetails }>}
   */
  static generateGrades(session, predicates = DefaultPredicates) {
    // Kumpulkan semua studentId unik
    const studentMap = new Map()
    session.components.forEach((comp) => {
      comp.scores.forEach((s) => {
        if (!studentMap.has(s.studentId)) {
          studentMap.set(s.studentId, { studentId: s.studentId, studentName: s.studentName })
        }
      })
    })

    const grades = []
    studentMap.forEach((studentInfo, studentId) => {
      // Ambil skor untuk semua komponen, urut sesuai komponen
      const scores = session.components.map((comp) => {
        const entry = comp.scores.find((s) => s.studentId === studentId)
        return entry || { score: null }
      })
      const finalScore = AssessmentEngine.calculateFinalScore(session.components, scores)
      const predicate = AssessmentEngine.getPredicate(finalScore, predicates)
      const componentDetails = session.components.map((comp, idx) => ({
        componentName: comp.name,
        score: scores[idx]?.score ?? null,
        maxScore: comp.maxScore,
        weight: comp.weight,
      }))
      grades.push({
        studentId,
        studentName: studentInfo.studentName,
        finalScore,
        predicate,
        componentDetails,
      })
    })
    return grades
  }

  /**
   * Hitung statistik kelas: rata-rata, tertinggi, terendah.
   * @param {Array<{ finalScore: number }>} grades
   * @returns {{ average: number, highest: number, lowest: number, count: number }}
   */
  static calculateClassStats(grades) {
    if (grades.length === 0) return { average: 0, highest: 0, lowest: 0, count: 0 }
    let sum = 0
    let highest = -Infinity
    let lowest = Infinity
    grades.forEach((g) => {
      sum += g.finalScore
      if (g.finalScore > highest) highest = g.finalScore
      if (g.finalScore < lowest) lowest = g.finalScore
    })
    return {
      average: Math.round((sum / grades.length) * 100) / 100,
      highest,
      lowest,
      count: grades.length,
    }
  }

  /**
   * Membuat objek ScoreEntry baru.
   * @param {{ componentId: number|string, studentId: number|string, studentName: string, score: number|null, notes?: string }} data
   * @returns {{ componentId, studentId, studentName, score, notes }}
   */
  static createScoreEntry(data) {
    return {
      componentId: data.componentId,
      studentId: data.studentId,
      studentName: data.studentName,
      score: data.score ?? null,
      notes: data.notes || '',
    }
  }

  /**
   * Mendapatkan daftar siswa yang belum dinilai pada komponen tertentu.
   * @param {Object} component - AssessmentComponent
   * @param {Array<{ studentId: number|string }>} allStudents
   * @returns {Array<{ studentId: number|string }>}
   */
  static getUnscoredStudents(component, allStudents) {
    const scoredIds = new Set(component.scores.map((s) => s.studentId))
    return allStudents.filter((s) => !scoredIds.has(s.studentId))
  }
}
