// src/models/AssessmentComponent.js
export class AssessmentComponent {
  constructor({
    id,
    sessionId,
    type, // 'PH', 'PTS', 'PAS', 'TUGAS', 'PRAKTIK', 'PROYEK', 'PORTOFOLIO', 'SIKAP'
    name, // misal: "PH 1 - Bab 1"
    weight, // bobot dalam persen (0-100)
    maxScore, // skor maksimal (misal 100)
    minScore = 0, // skor minimal (default 0)
    scores = [], // array of ScoreEntry
    isMandatory = true, // wajib diisi?
  }) {
    this.id = id
    this.sessionId = sessionId
    this.type = type
    this.name = name
    this.weight = weight
    this.maxScore = maxScore
    this.minScore = minScore
    this.scores = scores
    this.isMandatory = isMandatory
  }
}
