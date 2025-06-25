#!/bin/bash

# Quick rebuild - keeps database data
echo "🔄 Quick rebuilding backend..."
docker compose stop backend
docker compose build backend
docker compose up -d backend
echo "✅ Done! Check logs with: docker compose logs -f backend"