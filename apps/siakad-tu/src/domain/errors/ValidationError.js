import { ApplicationError } from './ApplicationError.js'

export class ValidationError extends ApplicationError {
  constructor(message = 'Validation failed', code = 'VALIDATION_ERROR', details = null) {
    super(message, code, 400, details)
  }
}
