import { BaseEntity } from '@/domain/shared/BaseEntity'

export class TeacherNote extends BaseEntity {
  constructor(props) {
    super(props)
    this.studentId = props.studentId
    this.teacherId = props.teacherId
    this.content = props.content || ''
    this.type = props.type || 'GENERAL' // 'GENERAL', 'ACHIEVEMENT', 'WARNING', 'COUNSELING'
  }

  isValid() {
    return this.content && this.content.trim().length > 0 && this.content.trim().length <= 500
  }

  toJSON() {
    return {
      ...super.toJSON(),
      studentId: this.studentId,
      teacherId: this.teacherId,
      content: this.content,
      type: this.type,
    }
  }
}
