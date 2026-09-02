export class StudentProgressEngine {
  /**
   * Hitung status kenaikan kelas
   * @param {number} averageScore - Rata-rata nilai semua mapel
   * @param {Object} coreSubjectsGrades - { subjectId: score } untuk mapel inti (misal: Matematika, B.Indonesia, dll)
   * @param {number} violationPoints - Total poin pelanggaran
   * @param {number} minScoreThreshold - Nilai minimal (default 60)
   * @param {number} maxViolationPoints - Maksimal poin pelanggaran (default 50)
   * @returns {Object} { status: 'NAIK'|'TIDAK_NAIK'|'PERHATIAN', reasons: [] }
   */
  static calculatePromotionStatus(
    averageScore,
    coreSubjectsGrades = {},
    violationPoints = 0,
    minScoreThreshold = 60,
    maxViolationPoints = 50,
  ) {
    const reasons = []

    // 1. Cek rata-rata
    if (averageScore < minScoreThreshold) {
      reasons.push(
        `Rata-rata nilai (${averageScore}) di bawah standar minimal (${minScoreThreshold})`,
      )
    }

    // 2. Cek nilai mapel inti (tidak boleh D)
    const coreSubjectsWithD = []
    for (const [subjectId, score] of Object.entries(coreSubjectsGrades)) {
      if (score < minScoreThreshold) {
        coreSubjectsWithD.push(subjectId)
      }
    }
    if (coreSubjectsWithD.length > 0) {
      reasons.push(`Terdapat ${coreSubjectsWithD.length} mata pelajaran inti di bawah KKM`)
    }

    // 3. Cek pelanggaran
    if (violationPoints >= maxViolationPoints) {
      reasons.push(
        `Total poin pelanggaran (${violationPoints}) melebihi batas maksimal (${maxViolationPoints})`,
      )
    }

    // 4. Tentukan status
    if (reasons.length === 0) {
      return { status: 'NAIK', reasons: ['Semua kriteria terpenuhi'] }
    }

    // Jika hanya pelanggaran yang menjadi masalah -> PERHATIAN (bisa naik dengan catatan)
    const isOnlyViolationIssue = reasons.length === 1 && reasons[0].includes('poin pelanggaran')
    if (isOnlyViolationIssue) {
      return { status: 'PERHATIAN', reasons }
    }

    return { status: 'TIDAK_NAIK', reasons }
  }

  /**
   * Hitung total poin prestasi dan pelanggaran
   */
  static calculateTotalPoints(achievements, violations) {
    const achievementPoints = achievements.reduce((sum, a) => sum + (a.points || 0), 0)
    const violationPoints = violations.reduce((sum, v) => sum + (v.points || 0), 0)
    return { achievementPoints, violationPoints }
  }

  /**
   * Berikan rekomendasi berdasarkan progress
   */
  static getRecommendation(progress) {
    const { averageScore, achievementPoints, violationPoints } = progress
    if (averageScore >= 85 && achievementPoints > 20) {
      return 'Sangat Baik (Pertahankan prestasi dan terus tingkatkan)'
    }
    if (averageScore < 60) {
      return 'Perlu Bimbingan Akademik (Rendahnya nilai akademik)'
    }
    if (violationPoints > 40) {
      return 'Perlu Bimbingan Perilaku (Poin pelanggaran tinggi)'
    }
    if (averageScore >= 70 && violationPoints <= 20) {
      return 'Baik (Terus tingkatkan)'
    }
    return 'Cukup (Perlu peningkatan di beberapa aspek)'
  }

  /**
   * Validasi catatan
   */
  static validateNote(content) {
    if (!content || content.trim().length === 0) {
      throw new Error('Catatan tidak boleh kosong.')
    }
    if (content.trim().length > 500) {
      throw new Error('Catatan maksimal 500 karakter.')
    }
    return true
  }

  /**
   * Validasi prestasi
   */
  static validateAchievement(achievement) {
    if (!achievement.title || achievement.title.trim() === '') {
      throw new Error('Judul prestasi tidak boleh kosong.')
    }
    if (achievement.points < 0 || achievement.points > 100) {
      throw new Error('Poin prestasi harus antara 0-100.')
    }
    return true
  }

  /**
   * Validasi pelanggaran
   */
  static validateViolation(violation) {
    if (!violation.title || violation.title.trim() === '') {
      throw new Error('Judul pelanggaran tidak boleh kosong.')
    }
    if (violation.points < 0 || violation.points > 100) {
      throw new Error('Poin pelanggaran harus antara 0-100.')
    }
    return true
  }
}
