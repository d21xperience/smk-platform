package repository

import (
    "context"

    "sekolah-platform/services/tu-core/internal/domain/inventory/aggregate"
)

type AssetIssueRepository interface {
    Save(ctx context.Context, issue *aggregate.AssetIssue) error
    FindByID(ctx context.Context, id string) (*aggregate.AssetIssue, error)
    FindByReporter(ctx context.Context, reporterID string) ([]*aggregate.AssetIssue, error)
    FindPending(ctx context.Context) ([]*aggregate.AssetIssue, error) // Reported + Verified
    FindByRoom(ctx context.Context, roomID string) ([]*aggregate.AssetIssue, error)
}
