import { ApplicationError } from './ApplicationError.js'

export class NotFoundError extends ApplicationError {
  constructor(message = 'Resource not found', code = 'NOT_FOUND', details = null) {
    super(message, code, 404, details)
  }
}
