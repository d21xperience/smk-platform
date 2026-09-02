// PURE JS - Value Object
export class Reflection {
  constructor(props) {
    this.studentUnderstanding = props.studentUnderstanding || ''
    this.teacherSelfReflection = props.teacherSelfReflection || ''
    this.improvementPlan = props.improvementPlan || ''
    this.notes = props.notes || ''
  }

  isEmpty() {
    return !this.studentUnderstanding && !this.teacherSelfReflection && !this.improvementPlan
  }

  toJSON() {
    return {
      studentUnderstanding: this.studentUnderstanding,
      teacherSelfReflection: this.teacherSelfReflection,
      improvementPlan: this.improvementPlan,
      notes: this.notes,
    }
  }

  static fromJSON(json) {
    return new Reflection(json)
  }
}
