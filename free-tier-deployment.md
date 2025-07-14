# AWS Free Tier Deployment Guide

This guide uses AWS Free Tier services to deploy your Go application for FREE testing.

## 🆓 AWS Free Tier Services We'll Use

- **EC2 (t2.micro)**: 750 hours/month FREE
- **RDS PostgreSQL**: 750 hours/month FREE (db.t3.micro)
- **Application Load Balancer**: FREE for 12 months
- **ECR**: FREE for 12 months (500MB storage)
- **CloudWatch**: FREE for 12 months

## 🚀 Quick Free Deployment Steps

### Step 1: Install AWS CLI and Docker
```bash
# Install AWS CLI (Windows)
winget install -e --id Amazon.AWSCLI

# Install Docker Desktop
# Download from: https://www.docker.com/products/docker-desktop/
```

### Step 2: Configure AWS CLI
```bash
aws configure
```
Enter your AWS credentials.

### Step 3: Test Locally (FREE)
```bash
# Test your app locally first
docker-compose up --build

# Test the API
curl http://localhost:8080/health
```

### Step 4: Simple EC2 Deployment (FREE)

Instead of complex ECS, let's use a simple EC2 instance:

```bash
# Create a simple deployment script
./simple-deploy.sh
```

## 📁 Files for Free Deployment

I'll create simplified deployment files that work with Free Tier:

1. **`simple-deploy.sh`** - Simple EC2 deployment
2. **`ec2-userdata.sh`** - EC2 setup script
3. **`free-tier-cloudformation.yml`** - Minimal infrastructure

## 💰 Total Cost: $0 (Free Tier)

- EC2 t2.micro: FREE (750 hours/month)
- RDS db.t3.micro: FREE (750 hours/month)
- Load Balancer: FREE (12 months)
- Data Transfer: FREE (15GB/month)

## ⚡ Quick Start Commands

```bash
# 1. Test locally
docker-compose up --build

# 2. Deploy to AWS (FREE)
./simple-deploy.sh

# 3. Get your app URL
aws ec2 describe-instances --filters "Name=tag:Name,Values=bharat-seva-app" --query "Reservations[].Instances[].PublicDnsName" --output text
```

## 🔧 What You'll Get

- **Your app running on**: `http://your-ec2-ip:8080`
- **Database**: PostgreSQL on RDS
- **Load Balancer**: For high availability
- **Monitoring**: CloudWatch logs

## 📊 Free Tier Limits

- **EC2**: 750 hours/month (enough for 24/7)
- **RDS**: 750 hours/month
- **Storage**: 20GB EBS storage
- **Data Transfer**: 15GB/month

## 🎯 Benefits of This Approach

1. **Completely FREE** for testing
2. **Simple setup** - no complex ECS
3. **Easy to understand** and modify
4. **Quick deployment** - 10-15 minutes
5. **Real production-like** environment

Would you like me to create the simplified deployment files now? 