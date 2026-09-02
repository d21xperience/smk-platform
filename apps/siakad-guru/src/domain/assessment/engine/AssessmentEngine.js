import { AssessmentSession } from '../models/AssessmentSession.js'
import { AssessmentComponent } from '../models/AssessmentComponent.js'
import {
  ASSESSMENT_STATUS,
  GRADE,
  PREDICATE,
  isValidAssessmentStatus,
  canFinalize,
} from '../models/AssessmentStatus.js'

export class AssessmentEngine {
  static validateSession(session) {
    if (!session || !(session instanceof AssessmentSession)) {
      throw new Error('Invalid assessment session object')
    }
    if (!session.isValid()) {
      throw new Error('Assessment session is missing required fields')
    }
    if (!isValidAssessmentStatus(session.status)) {
      throw new Error(`Invalid assessment status: ${session.status}`)
    }
    return true
  }

  static validateComponent(component) {
    if (!component || !(component instanceof AssessmentComponent)) {
      throw new Error('Invalid assessment component object')
    }
    if (!component.isValid()) {
      throw new Error('Assessment component is missing required fields or has invalid weight')
    }
    return true
  }

  static validateComponents(components) {
    if (!Array.isArray(components)) {
      throw new Error('Components must be an array')
    }
    if (components.length === 0) {
      throw new Error('At least one assessment component is required')
    }

    let totalWeight = 0
    components.forEach((component) => {
      AssessmentEngine.validateComponent(component)
      totalWeight += component.weight
    })

    if (totalWeight !== 100) {
      throw new Error(`Total component weight must be 100, got ${totalWeight}`)
    }

    return true
  }

  static calculateFinalScore(session) {
    AssessmentEngine.validateSession(session)

    if (session.components.length === 0) {
      throw new Error('No assessment components defined')
    }

    let weightedSum = 0
    let totalWeight = 0

    session.components.forEach((component) => {
      const score = session.scores[component.id]
      if (score !== undefined && score !== null) {
        if (score < 0 || score > 100) {
          throw new Error(`Invalid score for component ${component.name}: ${score}`)
        }
        weightedSum += score * component.weight
        totalWeight += component.weight
      }
    })

    if (totalWeight === 0) {
      return 0
    }

    return Math.round(weightedSum / totalWeight)
  }

  static calculateGrade(finalScore) {
    if (finalScore >= 85) return GRADE.A
    if (finalScore >= 70) return GRADE.B
    if (finalScore >= 55) return GRADE.C
    if (finalScore >= 40) return GRADE.D
    return GRADE.E
  }

  static calculatePredicate(finalScore) {
    if (finalScore >= 85) return PREDICATE.EXCELLENT
    if (finalScore >= 70) return PREDICATE.GOOD
    if (finalScore >= 55) return PREDICATE.SUFFICIENT
    if (finalScore >= 40) return PREDICATE.POOR
    return PREDICATE.VERY_POOR
  }

  static calculateResult(session) {
    const finalScore = AssessmentEngine.calculateFinalScore(session)
    const grade = AssessmentEngine.calculateGrade(finalScore)
    const predicate = AssessmentEngine.calculatePredicate(finalScore)

    return {
      finalScore,
      grade,
      predicate,
    }
  }

  static finalizeSession(session) {
    AssessmentEngine.validateSession(session)

    if (!canFinalize(session.status)) {
      throw new Error(`Cannot finalize session with status: ${session.status}`)
    }

    if (session.components.length === 0) {
      throw new Error('Cannot finalize assessment without components')
    }

    AssessmentEngine.validateComponents(session.components)

    session.status = ASSESSMENT_STATUS.FINALIZED
    return session
  }

  static updateScore(session, componentId, score) {
    AssessmentEngine.validateSession(session)

    if (session.isFinalized()) {
      throw new Error('Cannot update scores of a finalized assessment')
    }

    const component = session.components.find((c) => c.id === componentId)
    if (!component) {
      throw new Error(`Component not found: ${componentId}`)
    }

    if (score < 0 || score > 100) {
      throw new Error(`Score must be between 0 and 100, got ${score}`)
    }

    session.scores[componentId] = score
    return session
  }
}
