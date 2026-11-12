package service

import (
	"context"
	"sekolah-platform/services/tu-core/internal/repository"
)

type DashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetSummary(ctx context.Context, schoolID string, semesterID string) (map[string]interface{}, error) {
	kpiStats, _ := s.repo.GetKPIStats(ctx, schoolID, semesterID)
	jurusan, _ := s.repo.GetJurusanDistribusi(ctx, schoolID, semesterID)
	kehadiran, _ := s.repo.GetKehadiranHariIni(ctx, schoolID, semesterID)
	agenda, _ := s.repo.GetAgendaMendatang(ctx, schoolID, semesterID)

	ringkasanStaff := []map[string]interface{}{
		{"label": "Guru Tetap (PNS/Yayasan)", "jumlah": 54},
		{"label": "Guru Honorer", "jumlah": 22},
		{"label": "Tenaga Kependidikan", "jumlah": 8},
		{"label": "Wali Kelas Aktif", "jumlah": 42},
	}

	// Tambahkan variabel ini sebelum 'return'
	kpiList := []map[string]interface{}{
		{"label": "Total Siswa", "value": kpiStats["totalSiswa"], "icon": "groups", "color": "primary", "trend": 3.2, "routeTo": "/kesiswaan/siswa"},
		{"label": "Total Guru & Staf", "value": kpiStats["totalGuru"], "icon": "badge", "color": "secondary", "trend": 1.1, "routeTo": "/kepegawaian/data-guru"},
		{"label": "Total Rombel", "value": kpiStats["totalRombel"], "icon": "meeting_room", "color": "teal", "trend": 0},
		{"label": "Kompetensi Keahlian", "value": 6, "icon": "school", "color": "deep-orange", "trend": nil},
	}

	return map[string]interface{}{
		"kpis":              kpiList, // <-- Ganti map kpis yang lama dengan kpiList ini
		"totalSiswa":        kpiStats["totalSiswa"],
		"jurusanDistribusi": jurusan,
		"kehadiran":         kehadiran,
		"ringkasanStaff":    ringkasanStaff,
		"aktivitas":         []map[string]interface{}{{"id": 1, "judul": "Aktivitas sistem", "waktu": "Baru saja", "icon": "info", "warna": "primary"}},
		"agendaList":        agenda,
		"pengumumanList":    []map[string]interface{}{{"judul": "Pengumuman sistem", "tanggal": "Hari ini", "urgent": false}},
	}, nil
}
