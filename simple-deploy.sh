#!/bin/bash

# Simple AWS Free Tier Deployment Script
# This script deploys your Go app to EC2 for FREE testing

set -e

# Configuration
AWS_REGION="us-east-1"
STACK_NAME="bharat-seva-free"
EC2_INSTANCE_TYPE="t2.micro"
RDS_INSTANCE_CLASS="db.t3.micro"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 Starting FREE AWS Deployment for Bharat Seva Space${NC}"
echo -e "${YELLOW}This deployment uses AWS Free Tier services (Total Cost: $0)${NC}"

# Step 1: Create ECR Repository
echo -e "${YELLOW}Step 1: Creating ECR Repository...${NC}"
aws ecr create-repository --repository-name bharat-seva-space --region $AWS_REGION || echo "Repository already exists"

# Step 2: Login to ECR
echo -e "${YELLOW}Step 2: Logging into ECR...${NC}"
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com

# Step 3: Build and Push Docker Image
echo -e "${YELLOW}Step 3: Building and pushing Docker image...${NC}"
docker build -t bharat-seva-space .
ECR_URI=$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com/bharat-seva-space
docker tag bharat-seva-space:latest $ECR_URI:latest
docker push $ECR_URI:latest

# Step 4: Deploy CloudFormation Stack
echo -e "${YELLOW}Step 4: Deploying infrastructure (FREE tier)...${NC}"
aws cloudformation deploy \
  --template-file free-tier-cloudformation.yml \
  --stack-name $STACK_NAME \
  --parameter-overrides \
    DBPassword=YourSecurePassword123! \
    InstanceType=$EC2_INSTANCE_TYPE \
    DBInstanceClass=$RDS_INSTANCE_CLASS \
  --capabilities CAPABILITY_NAMED_IAM \
  --region $AWS_REGION

# Step 5: Get deployment outputs
echo -e "${YELLOW}Step 5: Getting deployment information...${NC}"
ALB_DNS=$(aws cloudformation describe-stacks \
  --stack-name $STACK_NAME \
  --query "Stacks[0].Outputs[?OutputKey=='LoadBalancerDNS'].OutputValue" \
  --output text \
  --region $AWS_REGION)

EC2_IP=$(aws cloudformation describe-stacks \
  --stack-name $STACK_NAME \
  --query "Stacks[0].Outputs[?OutputKey=='EC2PublicIP'].OutputValue" \
  --output text \
  --region $AWS_REGION)

RDS_ENDPOINT=$(aws cloudformation describe-stacks \
  --stack-name $STACK_NAME \
  --query "Stacks[0].Outputs[?OutputKey=='DatabaseEndpoint'].OutputValue" \
  --output text \
  --region $AWS_REGION)

echo -e "${GREEN}✅ Deployment Complete!${NC}"
echo -e "${YELLOW}Your application is now running:${NC}"
echo -e "🌐 Load Balancer URL: http://$ALB_DNS"
echo -e "🖥️  Direct EC2 URL: http://$EC2_IP:8080"
echo -e "🗄️  Database Endpoint: $RDS_ENDPOINT"
echo -e ""
echo -e "${GREEN}💰 Total Cost: $0 (AWS Free Tier)${NC}"
echo -e ""
echo -e "${YELLOW}Test your API:${NC}"
echo -e "curl http://$ALB_DNS/health"
echo -e ""
echo -e "${YELLOW}To clean up (when done testing):${NC}"
echo -e "aws cloudformation delete-stack --stack-name $STACK_NAME --region $AWS_REGION" 