// apps/siakad-tu/src/boot/studentEvents.js

import { studentEventListener } from '../stores/student/listeners/studentEventListener.js'

/**
 * Boot file untuk inisialisasi Student Event Listener.
 *
 * Dipanggil otomatis oleh Quasar saat app start.
 */
export default async () => {
  studentEventListener.init()
}
