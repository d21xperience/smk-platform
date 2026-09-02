// apps/siakad-tu/src/events/context/ContextEvents.js

import { DomainEvent } from '../DomainEvent.js';

/**
 * Dispatched saat user memilih atau berganti Operational Context.
 * Seluruh modul harus mendengarkan event ini untuk reset state mereka.
 */
export class OperationalContextChanged extends DomainEvent {
  constructor(context) {
    super(
      'OperationalContextChanged',
      context.schoolId,
      {
        schoolId: context.schoolId,
        schoolName: context.schoolName,
        academicYear: context.academicYear,
        semester: context.semester,
        periodId: context.periodId
      }
    );
  }
}

/**
 * Dispatched saat context di-clear (logout atau reset)
 */
export class OperationalContextCleared extends DomainEvent {
  constructor() {
    super('OperationalContextCleared', 'system', {});
  }
}