# !/bin/bash
# deploy/backup.sh

set -e

BACKUP_DIR="/backup/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

echo "💾 Starting backup at $(date)"

# Backup PostgreSQL
echo "🗄️ Backing up PostgreSQL..."
cd deploy/local
docker-compose exec -T postgres pg_dump -U $POSTGRES_USER -F c smk_platform > $BACKUP_DIR/postgres.dump

# Backup Redis (RDB file)
echo "📦 Backing up Redis..."
docker-compose exec -T redis redis-cli BGSAVE
sleep 5
docker cp $(docker-compose ps -q redis):/data/dump.rdb $BACKUP_DIR/redis.rdb

# Backup MinIO
echo "📁 Backing up MinIO..."
docker run --rm -v minio-data:/data -v $BACKUP_DIR:/backup alpine tar czf /backup/minio.tar.gz /data

# Backup SQLite (Edge)
echo "📄 Backing up SQLite..."
cd ../cloud
docker cp $(docker-compose ps -q Edge-bff):/data/Edge.db $BACKUP_DIR/Edge.db

# Compress backup
echo "🗜️ Compressing backup..."
cd /backup
tar czf $BACKUP_DIR.tar.gz $(basename $BACKUP_DIR)
rm -rf $BACKUP_DIR

# Upload to remote storage (optional)
echo "☁️ Uploading to remote storage..."
aws s3 cp $BACKUP_DIR.tar.gz s3://smk-backups/

# Cleanup old backups (keep last 7 days)
echo "🧹 Cleaning up old backups..."
find /backup -name ".tar.gz" -mtime +7 -delete

echo "✅ Backup completed at $(date)"
echo "📦 Backup file: $BACKUP_DIR.tar.gz"
