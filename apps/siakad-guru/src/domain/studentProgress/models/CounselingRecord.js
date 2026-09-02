import { BaseEntity } from '@/domain/shared/BaseEntity'

export class CounselingRecord extends BaseEntity {
  constructor(props) {
    super(props)
    this.studentId = props.studentId
    this.counselorId = props.counselorId
    this.topic = props.topic || ''
    this.content = props.content || ''
    this.followUp = props.followUp || ''
    this.date = props.date || new Date().toISOString().split('T')[0]
  }

  toJSON() {
    return {
      ...super.toJSON(),
      studentId: this.studentId,
      counselorId: this.counselorId,
      topic: this.topic,
      content: this.content,
      followUp: this.followUp,
      date: this.date,
    }
  }
}
