import { ASSESSMENT_STATUS } from '../../domain/assessment/models/AssessmentStatus.js'

// Mock assessments storage (keyed by assessmentId)
let mockAssessments = {}
let assessmentCounter = 1

export class AssessmentMockAdapter {
  async loadAssessmentByClass({ classId, subjectId, academicYearId, semesterId }) {
    await new Promise((resolve) => setTimeout(resolve, 400))

    for (const key in mockAssessments) {
      const assessment = mockAssessments[key]
      if (
        assessment.classId === classId &&
        assessment.subjectId === subjectId &&
        assessment.academicYearId === academicYearId &&
        assessment.semesterId === semesterId
      ) {
        return assessment
      }
    }

    return null
  }

  async loadAssessmentById({ assessmentId }) {
    await new Promise((resolve) => setTimeout(resolve, 400))
    return mockAssessments[assessmentId] || null
  }

  async createAssessment({
    classId,
    className,
    subjectId,
    subjectName,
    date,
    components,
    schoolId,
    academicYearId,
    semesterId,
  }) {
    await new Promise((resolve) => setTimeout(resolve, 400))

    const existing = await this.loadAssessmentByClass({
      classId,
      subjectId,
      academicYearId,
      semesterId,
    })
    if (existing) {
      throw new Error('Assessment already exists for this class and subject')
    }

    const newAssessment = {
      id: `ASM-${String(assessmentCounter++).padStart(3, '0')}`,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      components,
      scores: {},
      status: ASSESSMENT_STATUS.DRAFT,
      academicYearId,
      semesterId,
      schoolId,
    }

    mockAssessments[newAssessment.id] = newAssessment
    return newAssessment
  }

  async updateScore({ assessmentId, componentId, score }) {
    await new Promise((resolve) => setTimeout(resolve, 300))

    const assessment = mockAssessments[assessmentId]
    if (!assessment) {
      throw new Error(`Assessment not found: ${assessmentId}`)
    }

    if (assessment.status !== ASSESSMENT_STATUS.DRAFT) {
      throw new Error('Cannot update scores of a finalized assessment')
    }

    const component = assessment.components.find((c) => c.id === componentId)
    if (!component) {
      throw new Error(`Component not found: ${componentId}`)
    }

    if (score < 0 || score > 100) {
      throw new Error(`Score must be between 0 and 100, got ${score}`)
    }

    assessment.scores[componentId] = score
    return assessment
  }

  async finalizeAssessment({ assessmentId }) {
    await new Promise((resolve) => setTimeout(resolve, 400))

    const assessment = mockAssessments[assessmentId]
    if (!assessment) {
      throw new Error(`Assessment not found: ${assessmentId}`)
    }

    if (assessment.status !== ASSESSMENT_STATUS.DRAFT) {
      throw new Error('Assessment has already been finalized')
    }

    if (assessment.components.length === 0) {
      throw new Error('Cannot finalize assessment without components')
    }

    const totalWeight = assessment.components.reduce((sum, c) => sum + c.weight, 0)
    if (totalWeight !== 100) {
      throw new Error(`Total component weight must be 100, got ${totalWeight}`)
    }

    assessment.status = ASSESSMENT_STATUS.FINALIZED
    return assessment
  }
}
