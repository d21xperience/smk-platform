import { ApplicationError } from './ApplicationError.js';

export class ConflictError extends ApplicationError {
  constructor(message = 'Resource conflict', code = 'CONFLICT', details = null) {
    super(message, code, 409, details);
  }
}
