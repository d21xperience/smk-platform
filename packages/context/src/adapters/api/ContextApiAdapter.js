// === SET FALSE UNTUK MENGGUNAKAN API BACKEND ASLI ===
const MOCK_MODE = false;

const MOCK_PERIODS = []; // Tidak dipakai lagi

const DEFAULT_BASE_URL = "/api/v1";

export const createContextApiAdapter = (baseURL = DEFAULT_BASE_URL) => ({
  async getAvailablePeriods(context) {
    if (MOCK_MODE) return MOCK_PERIODS;

    // Menggunakan fetch. Pastikan endpoint ini tidak membutuhkan header auth
    // atau tangkap error-nya dengan baik.
    const res = await fetch(`${baseURL}/context/periods`, {
      headers: { "x-school-id": context.schoolId, "x-user-id": context.userId },
    });
    if (!res.ok) throw new Error("Failed to fetch periods");
    return await res.json();
  },

  async getUserPreference(context) {
    if (MOCK_MODE) return { semester_id: "20261" };

    const res = await fetch(`${baseURL}/users/me/preferences`, {
      headers: { "x-school-id": context.schoolId, "x-user-id": context.userId },
    });
    if (!res.ok) return null;
    return await res.json();
  },

  async saveUserPreference(context, semesterId) {
    if (MOCK_MODE) return { success: true, semester_id: semesterId };

    const res = await fetch(`${baseURL}/users/me/preferences`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "x-school-id": context.schoolId,
        "x-user-id": context.userId,
      },
      body: JSON.stringify({ semester_id: semesterId }),
    });
    if (!res.ok) throw new Error("Failed to save preference");
    return await res.json();
  },
});

export const ContextApiAdapter = createContextApiAdapter();
