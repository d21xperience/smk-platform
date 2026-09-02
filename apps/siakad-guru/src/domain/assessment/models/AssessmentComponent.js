export class AssessmentComponent {
  constructor({ id, name, type, weight }) {
    this.id = id
    this.name = name
    this.type = type
    this.weight = weight
  }

  isValid() {
    return !!(this.id && this.name && this.type && this.weight > 0 && this.weight <= 100)
  }
}
