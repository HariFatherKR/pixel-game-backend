#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔄 Starting Docker rebuild and restart process...${NC}"
echo ""

# Stop all containers
echo -e "${YELLOW}📦 Stopping containers...${NC}"
docker compose down
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Containers stopped successfully${NC}"
else
    echo -e "${RED}❌ Failed to stop containers${NC}"
    exit 1
fi

echo ""

# Remove old images (optional, uncomment if needed)
# echo -e "${YELLOW}🗑️  Removing old images...${NC}"
# docker compose rm -f

# Rebuild the backend image
echo -e "${YELLOW}🔨 Building backend image...${NC}"
docker compose build backend --no-cache
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Backend image built successfully${NC}"
else
    echo -e "${RED}❌ Failed to build backend image${NC}"
    exit 1
fi

echo ""

# Start all services
echo -e "${YELLOW}🚀 Starting all services...${NC}"
docker compose up -d
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ All services started successfully${NC}"
else
    echo -e "${RED}❌ Failed to start services${NC}"
    exit 1
fi

echo ""

# Wait for services to be healthy
echo -e "${YELLOW}⏳ Waiting for services to be healthy...${NC}"
sleep 5

# Check service health
echo -e "${YELLOW}🏥 Checking service health...${NC}"
docker compose ps

echo ""

# Test the API
echo -e "${YELLOW}🧪 Testing API health endpoint...${NC}"
sleep 2
curl -s http://localhost:8080/health > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ API is responding${NC}"
    echo ""
    echo -e "${GREEN}🎉 Rebuild complete!${NC}"
    echo ""
    echo -e "📝 Available endpoints:"
    echo -e "  - Health: http://localhost:8080/health"
    echo -e "  - Swagger: http://localhost:8080/swagger/index.html"
    echo -e "  - API: http://localhost:8080/api/v1/cards"
else
    echo -e "${RED}❌ API is not responding yet${NC}"
    echo -e "${YELLOW}Check logs with: docker compose logs backend${NC}"
fi

echo ""
echo -e "${YELLOW}📋 Useful commands:${NC}"
echo -e "  - View logs: docker compose logs -f backend"
echo -e "  - Stop services: docker compose down"
echo -e "  - Check status: docker compose ps"