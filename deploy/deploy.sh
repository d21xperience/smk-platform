#!/bin/bash
# deploy/deploy.sh

set -e  # Exit on error

ENVIRONMENT=${1:-cloud}  # cloud or local
VERSION=${2:-latest}

echo "🚀 Deploying to $ENVIRONMENT environment (version: $VERSION)"

# Load environment variables
if [ -f ".env.$ENVIRONMENT" ]; then
    export $(cat .env.$ENVIRONMENT | grep -v '^#' | xargs)
fi

# Function to deploy to Cloud VPS
deploy_cloud() {
    echo "📦 Deploying to Cloud VPS..."
    
    cd deploy/cloud
    
    # Pull latest images
    echo "📥 Pulling latest images..."
    docker-compose pull
    
    # Backup current state
    echo "💾 Creating backup..."
    docker-compose exec -T redis redis-cli BGSAVE
    sleep 5
    
    # Update frontend static files
    echo "🎨 Updating frontend..."
    docker-compose run --rm frontend-updater
    
    # Rolling update backend services
    echo "⚙️ Updating Edge BFF..."
    docker-compose up -d --no-deps edge-bff
    
    # Wait for health check
    echo "⏳ Waiting for Edge BFF to be healthy..."
    sleep 10
    
    # Verify deployment
    echo "✅ Verifying deployment..."
    if curl -f http://localhost/health; then
        echo "✅ Deployment successful!"
    else
        echo "❌ Deployment failed! Rolling back..."
        docker-compose rollback edge-bff
        exit 1
    fi
    
    # Cleanup old images
    echo "🧹 Cleaning up old images..."
    docker image prune -f
}

# Function to deploy to Local Server
deploy_local() {
    echo "📦 Deploying to Local Server..."
    
    cd deploy/local
    
    # Pull latest images
    echo "📥 Pulling latest images..."
    docker-compose pull
    
    # Run database migrations
    echo "🗄️ Running database migrations..."
    docker-compose run --rm migrator up
    
    # Backup database before update
    echo "💾 Creating database backup..."
    docker-compose exec -T postgres pg_dump -U $POSTGRES_USER smk_platform > backup/db_$(date +%Y%m%d_%H%M%S).sql
    
    # Rolling update core services
    echo "⚙️ Updating Core Services..."
    for service in core-academic core-student core-attendance core-assessment core-reporting; do
        echo "  Updating $service..."
        docker-compose up -d --no-deps $service
        sleep 5
    done
    
    # Verify deployment
    echo "✅ Verifying deployment..."
    for port in 50051 50052 50053 50054 50055; do
        if nc -z localhost $port; then
            echo "  ✅ Service on port $port is healthy"
        else
            echo "  ❌ Service on port $port is not responding!"
            exit 1
        fi
    done
    
    # Cleanup old images
    echo "🧹 Cleaning up old images..."
    docker image prune -f
    
    echo "✅ Local deployment successful!"
}

# Execute deployment
case $ENVIRONMENT in
    cloud)
        deploy_cloud
        ;;
    local)
        deploy_local
        ;;
    *)
        echo "❌ Unknown environment: $ENVIRONMENT"
        echo "Usage: $0 [cloud|local] [version]"
        exit 1
        ;;
esac