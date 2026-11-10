import { defineStore } from 'pinia';
import { ref } from 'vue';
import { api } from 'boot/axios';

export const useMaintenanceStore = defineStore('sarpras-maintenance', () => {
  const damageReports = ref([]);
  const schedules = ref([]);
  const logs = ref([]);
  const activeAssets = ref([]);
  const isLoading = ref(false);

  const fetchActiveAssets = async () => {
    try {
      const res = await api.get('/sarpras/kir/items');
      activeAssets.value = (res.data.data || []).filter(a => a.status_aset === 'AKTIF');
    } catch (e) {
      console.error('Gagal mengambil data aset:', e);
    }
  };

  // === DAMAGE REPORTS (TIKET KERUSAKAN) ===
  const fetchDamageReports = async () => {
    isLoading.value = true;
    try {
      const res = await api.get('/sarpras/damage-reports');
      damageReports.value = res.data.data || [];
    } catch (e) {
      console.error('Gagal mengambil data tiket:', e);
    } finally {
      isLoading.value = false;
    }
  };

  const createDamageReport = async (payload) => {
    const res = await api.post('/sarpras/damage-reports', payload);
    await fetchDamageReports();
    return res.data;
  };

  const updateDamageReport = async (id, payload) => {
    const res = await api.put(`/sarpras/damage-reports/${id}`, payload);
    await fetchDamageReports();
    return res.data;
  };

  // === MAINTENANCE SCHEDULES (JADWAL PERAWATAN) ===
  const fetchSchedules = async () => {
    isLoading.value = true;
    try {
      const res = await api.get('/sarpras/maintenance-schedules');
      schedules.value = res.data.data || [];
    } catch (e) {
      console.error('Gagal mengambil data jadwal:', e);
    } finally {
      isLoading.value = false;
    }
  };

  const createSchedule = async (payload) => {
    const res = await api.post('/sarpras/maintenance-schedules', payload);
    await fetchSchedules();
    return res.data;
  };

  const completeSchedule = async (id, payload) => {
    const res = await api.post(`/sarpras/maintenance-schedules/${id}/complete`, payload);
    await fetchSchedules();
    return res.data;
  };

  // === MAINTENANCE LOGS (BUKU PERAWATAN) ===
  const fetchMaintenanceLogs = async (startDate, endDate) => {
    isLoading.value = true;
    try {
      const res = await api.get(`/sarpras/maintenance-logs?start_date=${startDate || ''}&end_date=${endDate || ''}`);
      logs.value = res.data.data || [];
    } catch (e) {
      console.error('Gagal mengambil data log perawatan:', e);
    } finally {
      isLoading.value = false;
    }
  };

  const downloadMaintenanceBookPDF = async (startDate, endDate) => {
    const res = await api.get(`/sarpras/maintenance-logs?start_date=${startDate || ''}&end_date=${endDate || ''}&format=pdf`, {
      responseType: 'blob'
    });
    const url = window.URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', `Buku_Perawatan_${startDate || 'all'}_${endDate || 'all'}.pdf`);
    document.body.appendChild(link);
    link.click();
    link.remove();
  };

  return {
    damageReports, schedules, logs, activeAssets, isLoading,
    fetchActiveAssets, fetchDamageReports, createDamageReport, updateDamageReport,
    fetchSchedules, createSchedule, completeSchedule,
    fetchMaintenanceLogs, downloadMaintenanceBookPDF
  };
});
