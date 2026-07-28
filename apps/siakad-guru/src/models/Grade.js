// src/models/Grade.js
export class Grade {
  constructor({
    studentId,
    studentName,
    finalScore, // nilai akhir (0-100)
    predicate, // 'A', 'B', 'C', 'D' (dari Predicate)
    componentDetails, // array { componentName, weightedScore }
  }) {
    this.studentId = studentId
    this.studentName = studentName
    this.finalScore = finalScore
    this.predicate = predicate
    this.componentDetails = componentDetails
  }
}

// src/models/Predicate.js
export class Predicate {
  constructor({ minScore, maxScore, letter, description }) {
    this.minScore = minScore
    this.maxScore = maxScore
    this.letter = letter // 'A', 'B', 'C', 'D'
    this.description = description
  }
}
