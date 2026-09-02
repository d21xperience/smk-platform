import { BaseEntity } from '@/domain/shared/BaseEntity'

export class Achievement extends BaseEntity {
  constructor(props) {
    super(props)
    this.studentId = props.studentId
    this.title = props.title || ''
    this.description = props.description || ''
    this.category = props.category || 'ACADEMIC' // 'ACADEMIC', 'NON_ACADEMIC', 'CHARACTER'
    this.points = props.points || 0
    this.date = props.date || new Date().toISOString().split('T')[0]
  }

  isValid() {
    return this.title && this.title.trim() !== '' && this.points >= 0 && this.points <= 100
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
    }
  }
}
