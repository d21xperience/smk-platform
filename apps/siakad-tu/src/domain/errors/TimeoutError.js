import { ApplicationError } from './ApplicationError.js';

export class TimeoutError extends ApplicationError {
  constructor(message = 'Request timeout', code = 'TIMEOUT_ERROR', details = null) {
    super(message, code, 408, details);
  }
}
