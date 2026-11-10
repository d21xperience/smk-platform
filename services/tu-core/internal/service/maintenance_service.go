package service

import (
	"context"
	"errors"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
)

type MaintenanceService struct {
	repo repository.MaintenanceRepository
}

func NewMaintenanceService(repo repository.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{repo: repo}
}

func (s *MaintenanceService) CreateDamageReport(ctx context.Context, d *repository.DamageReport) (*repository.DamageReport, error) {
	if d.Description == "" || d.ReporterName == "" || d.AssetName == "" {
		return nil, errors.New("nama aset, pelapor, dan deskripsi wajib diisi")
	}
	reqNo, _ := s.repo.GenerateDamageNumber(ctx)
	d.ReportNumber = reqNo
	d.Status = "Pending"
	return d, s.repo.CreateDamageReport(ctx, d)
}

func (s *MaintenanceService) GetAllDamageReports(ctx context.Context, semesterID string) ([]*repository.DamageReport, error) {
	return s.repo.GetAllDamageReports(ctx, semesterID)
}

func (s *MaintenanceService) UpdateDamageReport(ctx context.Context, id string, status, vendor string, cost float64, notes string) error {
	if status == "" {
		return errors.New("status wajib diisi")
	}
	return s.repo.UpdateDamageReport(ctx, id, status, vendor, cost, notes)
}

func (s *MaintenanceService) CreateSchedule(ctx context.Context, sch *repository.MaintenanceSchedule) (*repository.MaintenanceSchedule, error) {
	if sch.AssetName == "" || sch.Frequency == "" {
		return nil, errors.New("nama aset dan frekuensi wajib diisi")
	}
	code, _ := s.repo.GenerateScheduleCode(ctx)
	sch.ScheduleCode = code
	sch.IsActive = true
	if sch.NextServiceDate.IsZero() {
		sch.NextServiceDate = time.Now()
	}
	return sch, s.repo.CreateSchedule(ctx, sch)
}

func (s *MaintenanceService) GetAllSchedules(ctx context.Context, semesterID string) ([]*repository.MaintenanceSchedule, error) {
	return s.repo.GetAllSchedules(ctx, semesterID)
}

func (s *MaintenanceService) CompleteSchedule(ctx context.Context, id string, serviceDate time.Time, cost float64, notes string, performedBy string) error {
	if serviceDate.IsZero() { serviceDate = time.Now() }
	return s.repo.CompleteSchedule(ctx, id, serviceDate, cost, notes, performedBy)
}
