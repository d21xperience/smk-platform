package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// BackupJob adalah job untuk backup database.
type BackupJob struct {
	dbHost     string
	dbPort     int
	dbName     string
	dbUser     string
	dbPassword string
	backupPath string
}

// NewBackupJob membuat BackupJob baru.
func NewBackupJob(dbHost string, dbPort int, dbName, dbUser, dbPassword, backupPath string) *BackupJob {
	return &BackupJob{
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbName:     dbName,
		dbUser:     dbUser,
		dbPassword: dbPassword,
		backupPath: backupPath,
	}
}

// Name mengimplementasikan scheduler.Job.
func (j *BackupJob) Name() string {
	return "database-backup"
}

// Execute mengimplementasikan scheduler.Job.
func (j *BackupJob) Execute(ctx context.Context) error {
	// Generate filename with timestamp
	filename := fmt.Sprintf("%s/backup_%s.sql", j.backupPath, time.Now().Format("20060102_150405"))

	// Run pg_dump
	cmd := exec.CommandContext(ctx, "pg_dump",
		"-h", j.dbHost,
		"-p", fmt.Sprintf("%d", j.dbPort),
		"-U", j.dbUser,
		"-d", j.dbName,
		"-F", "c", // Custom format
		"-f", filename,
	)

	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", j.dbPassword))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	return nil
}
