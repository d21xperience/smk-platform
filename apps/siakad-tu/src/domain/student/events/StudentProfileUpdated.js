import { DomainEvent } from '../../../events/DomainEvent.js'

// apps/siakad-tu/src/domain/student/events/StudentProfileUpdated.js
export class StudentProfileUpdated extends DomainEvent {
  constructor(studentId, schoolId, payload) {
    super('StudentProfileUpdated', studentId, { schoolId, ...payload })
  }
}
