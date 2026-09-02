import { BaseEntity } from '@/domain/shared/BaseEntity'

export class AssessmentFinalResult extends BaseEntity {
  constructor(props) {
    super(props)
    this.studentId = props.studentId
    this.finalScore = props.finalScore || 0
    this.details = props.details || {} // { componentId: score, ... }
    this.predicate = props.predicate || ''
    this.passed = props.passed || false
    this.calculatedAt = props.calculatedAt || new Date().toISOString()
  }

  update(finalScore, details) {
    this.finalScore = finalScore
    this.details = details || this.details
    this.calculatedAt = new Date().toISOString()
    this.updateTimestamps()
  }

  toJSON() {
    return {
      ...super.toJSON(),
      studentId: this.studentId,
      finalScore: this.finalScore,
      details: this.details,
      predicate: this.predicate,
      passed: this.passed,
      calculatedAt: this.calculatedAt,
    }
  }
}
