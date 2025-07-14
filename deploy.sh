#!/bin/bash

# AWS Deployment Script for Go Application
# Make sure you have AWS CLI configured and Docker installed

set -e

# Configuration
AWS_REGION="us-east-1"  # Change this to your preferred region
ECR_REPOSITORY_NAME="bharat-seva-space"
ECS_CLUSTER_NAME="bharat-seva-cluster"
ECS_SERVICE_NAME="bharat-seva-service"
ECS_TASK_DEFINITION_NAME="bharat-seva-task"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting AWS Deployment for Bharat Seva Space${NC}"

# Step 1: Create ECR Repository
echo -e "${YELLOW}Step 1: Creating ECR Repository...${NC}"
aws ecr create-repository --repository-name $ECR_REPOSITORY_NAME --region $AWS_REGION || echo "Repository already exists"

# Step 2: Get ECR Login Token
echo -e "${YELLOW}Step 2: Getting ECR Login Token...${NC}"
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com

# Step 3: Build and Tag Docker Image
echo -e "${YELLOW}Step 3: Building Docker Image...${NC}"
docker build -t $ECR_REPOSITORY_NAME .

# Step 4: Tag Image for ECR
echo -e "${YELLOW}Step 4: Tagging Image for ECR...${NC}"
ECR_URI=$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPOSITORY_NAME
docker tag $ECR_REPOSITORY_NAME:latest $ECR_URI:latest

# Step 5: Push Image to ECR
echo -e "${YELLOW}Step 5: Pushing Image to ECR...${NC}"
docker push $ECR_URI:latest

echo -e "${GREEN}✅ Docker image pushed successfully!${NC}"
echo -e "${YELLOW}Next steps:${NC}"
echo -e "1. Create RDS PostgreSQL instance"
echo -e "2. Create ECS Cluster"
echo -e "3. Create ECS Task Definition"
echo -e "4. Create ECS Service"
echo -e "5. Configure Application Load Balancer"
echo -e ""
echo -e "Run the following command to continue:"
echo -e "aws ecs create-cluster --cluster-name $ECS_CLUSTER_NAME --region $AWS_REGION" 