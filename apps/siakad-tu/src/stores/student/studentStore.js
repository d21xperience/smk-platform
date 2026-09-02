// apps/siakad-tu/src/stores/student/studentStore.js

import { useStudentListStore } from './studentListStore.js';
import { useStudentDetailStore } from './studentDetailStore.js';
import { useStudentFormStore } from './studentFormStore.js';
import { studentService } from '../../services/student/serviceFactory.js';

/**
 * studentStore — Facade yang menggabungkan 3 store dan menyediakan
 * high-level actions untuk Composable.
 *
 * Composable TIDAK perlu tahu tentang 3 store terpisah.
 * Cukup panggil method di facade ini.
 *
 * CATATAN: Facade ini TIDAK mengandung logika bisnis.
 * Dia hanya meneruskan panggilan ke Service dan mengupdate state.
 */
export function useStudentStore() {
  const listStore = useStudentListStore();
  const detailStore = useStudentDetailStore();
  const formStore = useStudentFormStore();

  /**
   * Fetch daftar siswa dengan filter & pagination saat ini
   * @param {OperationalContext} context
   */
  async function fetchList(context) {
    if (listStore.loading) return;

    listStore.setLoading(true);
    try {
      const result = await studentService.getStudents(context, {
        ...listStore.filters,
        page: listStore.page,
        limit: listStore.limit
      });

      if (result.success) {
        listStore.setItems(result.data);
      } else {
        listStore.setError(result.error);
      }
    } catch (error) {
      listStore.setError({ code: 'NETWORK_ERROR', message: error.message });
    }
  }

  /**
   * Fetch detail siswa by ID
   * @param {string} studentId
   * @param {OperationalContext} context
   * @param {boolean} forceRefresh - bypass cache
   */
  async function fetchDetail(studentId, context, forceRefresh = false) {
    // Cek cache dulu
    if (!forceRefresh) {
      const cached = detailStore.getById(studentId);
      if (cached) return cached;
    }

    if (detailStore.isLoading(studentId)) return null;

    detailStore.setLoading(studentId, true);
    try {
      const result = await studentService.getStudentById(studentId, context);

      if (result.success) {
        detailStore.setItem(studentId, result.data);
        return result.data;
      } else {
        detailStore.setError(studentId, result.error);
        return null;
      }
    } catch (error) {
      detailStore.setError(studentId, { code: 'NETWORK_ERROR', message: error.message });
      return null;
    }
  }

  /**
   * Register siswa baru
   * @param {Object} data
   * @param {OperationalContext} context
   */
  async function registerStudent(data, context) {
    formStore.setLoading(true);
    try {
      const result = await studentService.registerStudent(data, context);

      if (result.success) {
        listStore.addItem(result.data);
        detailStore.setItem(result.data.studentId, result.data);
        formStore.markSubmitted();
        return { success: true, data: result.data };
      } else {
        // Cek apakah error validasi (dari Engine) atau error lain
        if (result.error.code === 'VALIDATION_ERROR' || result.error.source === 'engine') {
          formStore.setValidationErrors([result.error]);
        } else {
          formStore.setError(result.error);
        }
        return { success: false, error: result.error };
      }
    } catch (error) {
      formStore.setError({ code: 'NETWORK_ERROR', message: error.message });
      return { success: false, error };
    }
  }

  /**
   * Enroll siswa ke kelas
   * @param {Object} data
   * @param {OperationalContext} context
   */
  async function enrollStudent(data, context) {
    listStore.setLoading(true);
    try {
      const result = await studentService.enrollStudent(data, context);

      if (result.success) {
        listStore.updateItem(result.data);
        detailStore.setItem(result.data.studentId, result.data);
        return { success: true, data: result.data };
      } else {
        listStore.setError(result.error);
        return { success: false, error: result.error };
      }
    } catch (error) {
      listStore.setError({ code: 'NETWORK_ERROR', message: error.message });
      return { success: false, error };
    } finally {
      listStore.setLoading(false);
    }
  }

  /**
   * Update profil siswa
   * @param {Object} data
   * @param {OperationalContext} context
   */
  async function updateStudentProfile(data, context) {
    formStore.setLoading(true);
    try {
      const result = await studentService.updateStudentProfile(data, context);

      if (result.success) {
        listStore.updateItem(result.data);
        detailStore.setItem(result.data.studentId, result.data);
        formStore.markSubmitted();
        return { success: true, data: result.data };
      } else {
        formStore.setError(result.error);
        return { success: false, error: result.error };
      }
    } catch (error) {
      formStore.setError({ code: 'NETWORK_ERROR', message: error.message });
      return { success: false, error };
    }
  }

  /**
   * Luluskan siswa
   * @param {string} studentId
   * @param {OperationalContext} context
   */
  async function graduateStudent(studentId, context) {
    listStore.setLoading(true);
    try {
      const result = await studentService.graduateStudent({ studentId }, context);

      if (result.success) {
        listStore.removeItem(studentId); // keluarkan dari list ACTIVE
        detailStore.setItem(studentId, result.data);
        return { success: true, data: result.data };
      } else {
        listStore.setError(result.error);
        return { success: false, error: result.error };
      }
    } catch (error) {
      listStore.setError({ code: 'NETWORK_ERROR', message: error.message });
      return { success: false, error };
    } finally {
      listStore.setLoading(false);
    }
  }

  /**
   * Pindahkan siswa
   * @param {Object} data
   * @param {OperationalContext} context
   */
  async function transferStudent(data, context) {
    listStore.setLoading(true);
    try {
      const result = await studentService.transferStudent(data, context);

      if (result.success) {
        listStore.removeItem(data.studentId);
        detailStore.setItem(data.studentId, result.data);
        return { success: true, data: result.data };
      } else {
        listStore.setError(result.error);
        return { success: false, error: result.error };
      }
    } catch (error) {
      listStore.setError({ code: 'NETWORK_ERROR', message: error.message });
      return { success: false, error };
    } finally {
      listStore.setLoading(false);
    }
  }

  /**
   * Reset semua store (saat context berubah)
   */
  function resetAll() {
    listStore.reset();
    detailStore.reset();
    formStore.reset();
  }

  return {
    // Sub-stores (untuk akses langsung jika perlu)
    list: listStore,
    detail: detailStore,
    form: formStore,

    // High-level actions
    fetchList,
    fetchDetail,
    registerStudent,
    enrollStudent,
    updateStudentProfile,
    graduateStudent,
    transferStudent,
    resetAll
  };
}
