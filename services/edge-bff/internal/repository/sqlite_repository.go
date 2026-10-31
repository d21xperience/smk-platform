// Package repository menyediakan SQLite read model untuk edge-bff.
package repository

import (
	"context"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// StudentProjection adalah read model untuk siswa di SQLite.
type StudentProjection struct {
	StudentID   string `gorm:"primaryKey"`
	SchoolID    string `gorm:"index"`
	NISN        string `gorm:"index"`
	NIS         string `gorm:"index"`
	FirstName   string
	MiddleName  string
	LastName    string
	BirthDate   time.Time
	Gender      string
	AddressCity string
	Status      string `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastEventID string // Untuk idempotency
}

func (StudentProjection) TableName() string { return "students" }

// TeacherProjection adalah read model untuk guru di SQLite.
type TeacherProjection struct {
	TeacherID   string `gorm:"primaryKey"`
	SchoolID    string `gorm:"index"`
	NIP         string `gorm:"index"`
	FirstName   string
	LastName    string
	Subject     string `gorm:"index"`
	Gender      string
	Status      string `gorm:"index"`
	UpdatedAt   time.Time
	LastEventID string
}

func (TeacherProjection) TableName() string { return "teachers" }

// ClassProjection adalah read model untuk kelas di SQLite.
type ClassProjection struct {
	ClassID     string `gorm:"primaryKey"`
	SchoolID    string `gorm:"index"`
	PeriodID    string `gorm:"index"`
	Name        string
	Grade       int
	Level       string
	Major       string
	HomeroomID  string
	Capacity    int
	UpdatedAt   time.Time
	LastEventID string
}

func (ClassProjection) TableName() string { return "classes" }

// SQLiteRepository adalah repository untuk SQLite read model.
type SQLiteRepository struct {
	db *gorm.DB
}

// NewSQLiteRepository membuat SQLiteRepository baru.
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	logLevel := logger.Silent
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&StudentProjection{},
		&TeacherProjection{},
		&ClassProjection{},
	); err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}

// DB mengembalikan underlying GORM DB (untuk advanced queries).
func (r *SQLiteRepository) DB() *gorm.DB {
	return r.db
}

// Close menutup koneksi.
func (r *SQLiteRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// === Student Operations ===

// UpsertStudent menginsert atau update student projection.
func (r *SQLiteRepository) UpsertStudent(ctx context.Context, p *StudentProjection) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// FindStudentByID mencari student by ID.
// Fungsi lama (JANGAN DIUBAH, dipakai oleh Projection Manager)
func (r *SQLiteRepository) FindStudentByID(ctx context.Context, studentID string) (*StudentProjection, error) {
	var student StudentProjection
	err := r.db.WithContext(ctx).Where("id = ?", studentID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *SQLiteRepository) FindTeacherByID(ctx context.Context, teacherID string) (*TeacherProjection, error) {
	var teacher TeacherProjection
	err := r.db.WithContext(ctx).Where("id = ?", teacherID).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *SQLiteRepository) UpdateStudentStatus(ctx context.Context, studentID, status, eventID string) error {
	// Sesuaikan isi query ini dengan aslinya, yang penting tidak ada schoolID
	return r.db.WithContext(ctx).Model(&StudentProjection{}).Where("id = ?", studentID).Updates(map[string]interface{}{"status": status, "event_id": eventID}).Error
}

// ListStudents mengembalikan daftar siswa dengan filter.
func (r *SQLiteRepository) ListStudents(ctx context.Context, schoolID, periodID, status, search string, page, limit int) ([]StudentProjection, int64, error) {
	query := r.db.WithContext(ctx).Model(&StudentProjection{}).Where("school_id = ?", schoolID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR nisn LIKE ? OR nis LIKE ?",
			like, like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var results []StudentProjection
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("first_name ASC").Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// === Teacher Operations ===

// UpsertTeacher menginsert atau update teacher projection.
func (r *SQLiteRepository) UpsertTeacher(ctx context.Context, p *TeacherProjection) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// ListTeachers mengembalikan daftar guru dengan filter.
func (r *SQLiteRepository) ListTeachers(ctx context.Context, schoolID, subject, status string, page, limit int) ([]TeacherProjection, int64, error) {
	query := r.db.WithContext(ctx).Model(&TeacherProjection{}).Where("school_id = ?", schoolID)

	if subject != "" {
		query = query.Where("subject = ?", subject)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var results []TeacherProjection
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("first_name ASC").Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// === Class Operations ===

// UpsertClass menginsert atau update class projection.
func (r *SQLiteRepository) UpsertClass(ctx context.Context, p *ClassProjection) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// ListClassesByPeriod mengembalikan daftar kelas untuk periode tertentu.
func (r *SQLiteRepository) ListClassesByPeriod(ctx context.Context, schoolID, periodID string) ([]ClassProjection, error) {
	var results []ClassProjection
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND period_id = ?", schoolID, periodID).
		Order("grade ASC, name ASC").
		Find(&results).Error
	return results, err
}

// === Idempotency ===

// HasProcessedEvent mengecek apakah event sudah diproses.
func (r *SQLiteRepository) HasProcessedEvent(ctx context.Context, eventID string) (bool, error) {
	// Check di semua projection tables
	var count int64
	err := r.db.WithContext(ctx).
		Model(&StudentProjection{}).
		Where("last_event_id = ?", eventID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	err = r.db.WithContext(ctx).
		Model(&TeacherProjection{}).
		Where("last_event_id = ?", eventID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindStudentByIDAndSchoolID mencari student berdasarkan ID dan SchoolID (untuk handler HTTP).
func (r *SQLiteRepository) FindStudentByIDAndSchoolID(ctx context.Context, studentID string, schoolID string) (*StudentProjection, error) {
	var student StudentProjection
	err := r.db.WithContext(ctx).Where("id = ? AND school_id = ?", studentID, schoolID).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// FindTeacherByIDAndSchoolID mencari teacher berdasarkan ID dan SchoolID (untuk handler HTTP).
func (r *SQLiteRepository) FindTeacherByIDAndSchoolID(ctx context.Context, teacherID string, schoolID string) (*TeacherProjection, error) {
	var teacher TeacherProjection
	err := r.db.WithContext(ctx).Where("id = ? AND school_id = ?", teacherID, schoolID).First(&teacher).Error
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// UpdateStudentStatusBySchoolID mengupdate status siswa berdasarkan ID dan SchoolID (untuk handler HTTP jika ada).
func (r *SQLiteRepository) UpdateStudentStatusBySchoolID(ctx context.Context, studentID string, schoolID string, status string, eventID string) error {
	result := r.db.WithContext(ctx).Model(&StudentProjection{}).
		Where("id = ? AND school_id = ?", studentID, schoolID).
		Updates(map[string]interface{}{
			"status":   status,
			"event_id": eventID,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
