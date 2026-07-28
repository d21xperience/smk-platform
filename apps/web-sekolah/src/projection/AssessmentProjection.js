import { DomainEventType } from 'src/events/DomainEventType'

export function assessmentProjection(state, event) {
  switch (event.type) {
    case DomainEventType.ASSESSMENT_SUBMITTED: {
      const { className, subject, stats } = event.payload
      const key = `${className}_${subject}`
      return {
        ...state,
        assessment: {
          ...state.assessment,
          [key]: {
            className,
            subject,
            average: stats.average,
            highest: stats.highest,
            lowest: stats.lowest,
            count: stats.count,
            submittedAt: event.timestamp,
          },
        },
      }
    }
    case DomainEventType.ASSESSMENT_DRAFT_SAVED: {
      // Bisa menambah indikator "draft tersimpan" jika diperlukan
      return state
    }
    default:
      return state
  }
}
