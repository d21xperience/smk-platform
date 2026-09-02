import { DomainEvent } from '../../../events/DomainEvent.js'

// apps/siakad-tu/src/domain/student/events/StudentEnrolled.js
export class StudentEnrolled extends DomainEvent {
  constructor(studentId, schoolId, payload) {
    super('StudentEnrolled', studentId, { schoolId, ...payload })
  }
}
