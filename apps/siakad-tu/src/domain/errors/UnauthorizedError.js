import { ApplicationError } from './ApplicationError.js';

export class UnauthorizedError extends ApplicationError {
  constructor(message = 'Unauthorized access', code = 'UNAUTHORIZED', details = null) {
    super(message, code, 401, details);
  }
}
