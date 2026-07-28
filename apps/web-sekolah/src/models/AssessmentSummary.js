export class AssessmentSummary {
  constructor({
    subject,
    finalScore,
    predicate,
    componentDetails = [], // [{ componentName, score, maxScore, weight }]
  }) {
    this.subject = subject
    this.finalScore = finalScore
    this.predicate = predicate
    this.componentDetails = componentDetails
  }
}
