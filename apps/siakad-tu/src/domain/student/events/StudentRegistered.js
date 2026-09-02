// apps/siakad-tu/src/domain/student/events/StudentRegistered.js
import { DomainEvent } from '../../../events/DomainEvent.js'

export class StudentRegistered extends DomainEvent {
  constructor(studentId, schoolId, payload) {
    super('StudentRegistered', studentId, { schoolId, ...payload })
  }
}


