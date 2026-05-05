#!/bin/bash

# ==============================================================================
# HEALY Monorepo - Automated Deployment Script (Standard IT Practice)
# ==============================================================================
# Usage: chmod +x deploy.sh && ./deploy.sh
# ==============================================================================

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Starting Deployment Process...${NC}"

# 1. Check for .env files
echo -e "${YELLOW}🔍 Checking Environment Files...${NC}"
if [ ! -f "backend/.env" ]; then
    echo -e "${RED}❌ Error: backend/.env not found!${NC}"
    echo -e "Please run: cp backend/.env.example backend/.env and configure it."
    exit 1
fi

if [ ! -f "frontend/.env.local" ]; then
    echo -e "${RED}❌ Error: frontend/.env.local not found!${NC}"
    echo -e "Please run: cp frontend/.env.example frontend/.env.local and configure it."
    exit 1
fi

# 2. Pull Latest Changes (Optional)
# echo -e "${YELLOW}📥 Pulling latest changes from Git...${NC}"
# git pull origin main

# 3. Backend Setup (Go)
echo -e "${YELLOW}🔧 Setting up Backend (Go)...${NC}"
cd backend
go mod tidy
echo -e "${YELLOW}🏗️ Building Backend Binary...${NC}"
go build -o healy-server ./cmd/api

# Restart Backend Process
echo -e "${YELLOW}🔄 Restarting Backend with PM2...${NC}"
pm2 delete healy-backend 2>/dev/null || true
pm2 start ./healy-server --name "healy-backend"
cd ..

# 4. Frontend Setup (Next.js)
echo -e "${YELLOW}🔧 Setting up Frontend (Next.js)...${NC}"
cd frontend
echo -e "${YELLOW}📦 Installing NPM Dependencies...${NC}"
npm install --frozen-lockfile || npm install
echo -e "${YELLOW}🏗️ Building Frontend (Next.js)...${NC}"
npm run build

# Restart Frontend Process
echo -e "${YELLOW}🔄 Restarting Frontend with PM2...${NC}"
pm2 delete healy-frontend 2>/dev/null || true
pm2 start npm --name "healy-frontend" -- start
cd ..

# 5. Final Status
echo -e "${GREEN}✅ Deployment Completed Successfully!${NC}"
echo -e "${YELLOW}📊 Current PM2 Status:${NC}"
pm2 status

echo -e "\n${GREEN}Tips: Use 'pm2 logs' to monitor your application.${NC}"
