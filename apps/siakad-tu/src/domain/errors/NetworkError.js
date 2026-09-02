import { ApplicationError } from './ApplicationError.js';

export class NetworkError extends ApplicationError {
  constructor(message = 'Network error occurred', code = 'NETWORK_ERROR', details = null) {
    super(message, code, 0, details);
  }
}
