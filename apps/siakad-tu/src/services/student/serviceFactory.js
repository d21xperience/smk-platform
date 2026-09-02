// apps/siakad-tu/src/services/student/serviceFactory.js

import { StudentEngine } from '@/engine/student/StudentEngine';
import { studentApi } from '@/adapters/api/studentApi';
import { eventDispatcher } from '@/events/eventDispatcher';
import { StudentService } from './StudentService.js';

/**
 * Service Factory — Dependency Injection container.
 *
 * Membuat instance StudentService dengan dependencies yang sudah di-inject.
 *
 * Keuntungan DI:
 * 1. Mudah di-test (bisa inject mock dependencies)
 * 2. Mudah di-swap (mock → real adapter)
 * 3. Singleton per application lifecycle
 * 4. Dependencies eksplisit, tidak tersembunyi
 *
 * Penggunaan:
 *   import { studentService } from './serviceFactory.js';
 *   await studentService.registerStudent(data, context);
 *
 * Untuk testing:
 *   import { createStudentService } from './serviceFactory.js';
 *   const testService = createStudentService({
 *     engine: mockEngine,
 *     adapter: mockAdapter,
 *     eventDispatcher: mockDispatcher
 *   });
 */

/**
 * Buat instance StudentService dengan dependencies custom
 * (untuk testing atau environment khusus)
 */
export function createStudentService(dependencies = {}) {
  const engine = dependencies.engine || new StudentEngine();
  const adapter = dependencies.adapter || studentApi;
  const eventDispatcherInstance = dependencies.eventDispatcher || eventDispatcher;

  return new StudentService({
    engine,
    adapter,
    eventDispatcher: eventDispatcherInstance
  });
}

/**
 * Singleton instance untuk penggunaan normal di aplikasi.
 * Dibuat sekali, dipakai di seluruh aplikasi.
 */
export const studentService = createStudentService();
