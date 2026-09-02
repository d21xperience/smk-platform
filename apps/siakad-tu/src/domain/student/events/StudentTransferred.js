import { DomainEvent } from '../../../events/DomainEvent.js'

// apps/siakad-tu/src/domain/student/events/StudentTransferred.js
export class StudentTransferred extends DomainEvent {
  constructor(studentId, schoolId, payload) {
    super('StudentTransferred', studentId, { schoolId, ...payload })
  }
}

