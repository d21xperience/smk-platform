import { ApplicationError } from './ApplicationError.js';

export class BusinessRuleError extends ApplicationError {
  constructor(message = 'Business rule violation', code = 'BUSINESS_RULE_ERROR', details = null) {
    super(message, code, 422, details);
  }
}
