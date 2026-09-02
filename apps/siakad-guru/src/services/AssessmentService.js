import { AssessmentEngine } from '../domain/assessment/engine/AssessmentEngine.js'
import { AssessmentSession } from '../domain/assessment/models/AssessmentSession.js'
import { AssessmentComponent } from '../domain/assessment/models/AssessmentComponent.js'
import { AssessmentFinalizedEvent } from '../domain/assessment/events/AssessmentFinalized.js'

export class AssessmentService {
  constructor({ assessmentAdapter, eventDispatcher }) {
    this.assessmentAdapter = assessmentAdapter
    this.eventDispatcher = eventDispatcher
  }

  async loadAssessment({ classId, subjectId, academicYearId, semesterId }) {
    const rawData = await this.assessmentAdapter.loadAssessmentByClass({
      classId,
      subjectId,
      academicYearId,
      semesterId,
    })

    if (!rawData) {
      return null
    }

    const components = rawData.components.map((c) => new AssessmentComponent(c))
    const session = new AssessmentSession({
      ...rawData,
      components,
    })

    AssessmentEngine.validateSession(session)
    return session
  }

  async loadAssessmentById({ assessmentId }) {
    const rawData = await this.assessmentAdapter.loadAssessmentById({ assessmentId })

    if (!rawData) {
      return null
    }

    const components = rawData.components.map((c) => new AssessmentComponent(c))
    const session = new AssessmentSession({
      ...rawData,
      components,
    })

    AssessmentEngine.validateSession(session)
    return session
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
    const domainComponents = components.map((c) => new AssessmentComponent(c))
    AssessmentEngine.validateComponents(domainComponents)

    const rawData = await this.assessmentAdapter.createAssessment({
      classId,
      className,
      subjectId,
      subjectName,
      date,
      components: components.map((c) => ({
        id: c.id,
        name: c.name,
        type: c.type,
        weight: c.weight,
      })),
      schoolId,
      academicYearId,
      semesterId,
    })

    const resultComponents = rawData.components.map((c) => new AssessmentComponent(c))
    const session = new AssessmentSession({
      ...rawData,
      components: resultComponents,
    })

    AssessmentEngine.validateSession(session)
    return session
  }

  async updateScore({ assessmentId, componentId, score }) {
    const rawData = await this.assessmentAdapter.updateScore({
      assessmentId,
      componentId,
      score,
    })

    const components = rawData.components.map((c) => new AssessmentComponent(c))
    const session = new AssessmentSession({
      ...rawData,
      components,
    })

    AssessmentEngine.validateSession(session)
    return session
  }

  async finalizeAssessment({ assessmentId }) {
    const rawData = await this.assessmentAdapter.finalizeAssessment({ assessmentId })

    const components = rawData.components.map((c) => new AssessmentComponent(c))
    const session = new AssessmentSession({
      ...rawData,
      components,
    })

    AssessmentEngine.validateSession(session)
    AssessmentEngine.finalizeSession(session)

    const result = AssessmentEngine.calculateResult(session)

    if (this.eventDispatcher) {
      this.eventDispatcher.dispatch(new AssessmentFinalizedEvent(session, result))
    }

    return { session, result }
  }

  calculatePreview(session) {
    return AssessmentEngine.calculateResult(session)
  }
}
