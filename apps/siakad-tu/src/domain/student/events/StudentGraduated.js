import { DomainEvent } from '../../../events/DomainEvent.js'

// apps/siakad-tu/src/domain/student/events/StudentGraduated.js
export class StudentGraduated extends DomainEvent {
  constructor(studentId, schoolId, payload) {
    super('StudentGraduated', studentId, { schoolId, ...payload })
  }
}

