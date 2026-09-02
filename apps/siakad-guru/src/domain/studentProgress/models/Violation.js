import { BaseEntity } from '@/domain/shared/BaseEntity'

export class Violation extends BaseEntity {
  constructor(props) {
    super(props)
    this.studentId = props.studentId
    this.title = props.title || ''
    this.description = props.description || ''
    this.category = props.category || 'DISCIPLINE' // 'DISCIPLINE', 'ACADEMIC', 'ATTENDANCE'
    this.points = props.points || 0
    this.date = props.date || new Date().toISOString().split('T')[0]
    this.isResolved = props.isResolved || false
  }

  isValid() {
    return this.title && this.title.trim() !== '' && this.points >= 0 && this.points <= 100
  }

  resolve() {
    this.isResolved = true
    this.updateTimestamps()
  }

  toJSON() {
    return {
      ...super.toJSON(),
      studentId: this.studentId,
      title: this.title,
      description: this.description,
      category: this.category,
      points: this.points,
      date: this.date,
      isResolved: this.isResolved,
    }
  }
}
