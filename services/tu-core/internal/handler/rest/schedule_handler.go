package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleHandler struct {
	classRepo    repository.ClassRepository
	academicRepo repository.AcademicRepository
	db           *pgxpool.Pool
}

func NewScheduleHandler(classRepo repository.ClassRepository, academicRepo repository.AcademicRepository, db *pgxpool.Pool) *ScheduleHandler {
	return &ScheduleHandler{classRepo: classRepo, academicRepo: academicRepo, db: db}
}

// GetOptions (GET /api/v1/schedules/options)
func (h *ScheduleHandler) GetOptions(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	// 1. Ambil List Kelas
	classes, err := h.academicRepo.ListClasses(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal memuat kelas: " + err.Error()})
		return
	}

	// 2. Ambil List Mapel (Gunakan COALESCE untuk jaga-jaga jika NULL)
	rowsMapel, err := h.db.Query(req.Context(), `SELECT mata_pelajaran_id, COALESCE(nm_mapel, 'Tanpa Nama') FROM tabel_mapel ORDER BY nm_mapel ASC`)
	var mapels []map[string]interface{}
	if err == nil {
		defer rowsMapel.Close()
		for rowsMapel.Next() {
			var id int
			var name string
			if err := rowsMapel.Scan(&id, &name); err == nil {
				mapels = append(mapels, map[string]interface{}{"id": id, "name": name})
			}
		}
	}

	// 3. Ambil List Guru
	rowsPtk, err := h.db.Query(req.Context(), `SELECT ptk_id::text, COALESCE(nama, 'Tanpa Nama') FROM tabel_ptk ORDER BY nama ASC`)
	var teachers []map[string]interface{}
	if err == nil {
		defer rowsPtk.Close()
		for rowsPtk.Next() {
			var id, name string
			if err := rowsPtk.Scan(&id, &name); err == nil {
				teachers = append(teachers, map[string]interface{}{"id": id, "name": name})
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"classes":  classes,
		"subjects": mapels,
		"teachers": teachers,
	})
}

// CreateSchedule (POST /api/v1/schedules)
func (h *ScheduleHandler) CreateSchedule(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	var body struct {
		ClassID   string `json:"class_id"`
		SubjectID int    `json:"subject_id"`
		TeacherID string `json:"teacher_id"`
		Day       string `json:"day"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Format JSON tidak valid"})
		return
	}

	// 1. Cari atau Buat Pembelajaran ID
	var pembelajaranID string
	queryFind := `SELECT pembelajaran_id FROM tabel_pembelajaran WHERE rombongan_belajar_id = $1 AND mata_pelajaran_id = $2 AND ptk_terdaftar_id = $3 AND semester_id = $4`
	err := h.db.QueryRow(req.Context(), queryFind, body.ClassID, body.SubjectID, body.TeacherID, opCtx.AcademicPeriodID).Scan(&pembelajaranID)

	if err != nil {
		// Jika belum ada, buat baru
		pembelajaranID = uuid.NewString()
		queryInsert := `INSERT INTO tabel_pembelajaran (pembelajaran_id, rombongan_belajar_id, mata_pelajaran_id, semester_id, ptk_terdaftar_id) VALUES ($1, $2, $3, $4, $5)`
		_, errInsert := h.db.Exec(req.Context(), queryInsert, pembelajaranID, body.ClassID, body.SubjectID, opCtx.AcademicPeriodID, body.TeacherID)
		if errInsert != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal membuat pembelajaran: " + errInsert.Error()})
			return
		}
	}

	// 2. Simpan ke tabel_jadwal_pelajaran
	scheduleID := uuid.NewString()
	querySchedule := `INSERT INTO tabel_jadwal_pelajaran (id, pembelajaran_id, hari, jam_mulai, jam_selesai) VALUES ($1, $2, $3, $4, $5)`
	_, err = h.db.Exec(req.Context(), querySchedule, scheduleID, pembelajaranID, body.Day, body.StartTime, body.EndTime)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Gagal menyimpan jadwal: " + err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "Jadwal berhasil ditambahkan"})
}

// ListSchedules (GET /api/v1/schedules)
func (h *ScheduleHandler) ListSchedules(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing context"})
		return
	}

	// PERBAIKAN: Gunakan COALESCE untuk semua kolom string agar tidak panic jika ada data NULL
	query := `
		SELECT
			j.id::text,
			COALESCE(c.nm_kelas, 'Tanpa Kelas'),
			COALESCE(m.nm_mapel, 'Tanpa Mapel'),
			COALESCE(p.nama, 'Belum Ada Guru'),
			j.hari,
			j.jam_mulai::text,
			j.jam_selesai::text
		FROM tabel_jadwal_pelajaran j
		JOIN tabel_pembelajaran tp ON j.pembelajaran_id = tp.pembelajaran_id
		JOIN tabel_kelas c ON tp.rombongan_belajar_id = c.rombongan_belajar_id
		JOIN tabel_mapel m ON tp.mata_pelajaran_id = m.mata_pelajaran_id
		LEFT JOIN tabel_ptk p ON tp.ptk_terdaftar_id = p.ptk_id
		WHERE c.sekolah_id = $1 AND tp.semester_id = $2
		ORDER BY j.hari ASC, j.jam_mulai ASC
	`
	rows, err := h.db.Query(req.Context(), query, opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		// Jika tabel_jadwal_pelajaran belum ada, error akan muncul di sini
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Query Error: " + err.Error()})
		return
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var id, kelas, mapel, guru, hari, jamMulai, jamSelesai string
		if err := rows.Scan(&id, &kelas, &mapel, &guru, &hari, &jamMulai, &jamSelesai); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Scan Error: " + err.Error()})
			return
		}
		result = append(result, map[string]interface{}{
			"id": id, "class_name": kelas, "subject_name": mapel, "teacher_name": guru,
			"day": hari, "start_time": jamMulai, "end_time": jamSelesai,
		})
	}

	if result == nil {
		result = []map[string]interface{}{}
	}
	writeJSON(w, http.StatusOK, result)
}
