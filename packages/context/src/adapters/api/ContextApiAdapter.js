// Default fallback jika tidak di-inject
const DEFAULT_BASE_URL = "/api/v1"; 

export const createContextApiAdapter = (baseURL = DEFAULT_BASE_URL) => ({
  async getAvailablePeriods(context) {
    const res = await fetch(`${baseURL}/context/periods`, {
      headers: { "x-school-id": context.schoolId, "x-user-id": context.userId },
    });
    if (!res.ok) throw new Error("Failed to fetch periods");
    return await res.json();
  },
  async getUserPreference(context) {
    const res = await fetch(`${baseURL}/users/me/preferences`, {
      headers: { "x-school-id": context.schoolId, "x-user-id": context.userId },
    });
    if (!res.ok) return null;
    return await res.json();
  },
  async saveUserPreference(context, semesterId) {
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

// Instance default
export const ContextApiAdapter = createContextApiAdapter();
