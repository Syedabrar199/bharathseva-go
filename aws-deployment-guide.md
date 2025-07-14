# AWS Deployment Guide for Bharat Seva Space

This guide will help you deploy your Go application to AWS using ECS (Elastic Container Service) with RDS PostgreSQL.

## Prerequisites

1. **AWS Account** ✅ (You have this)
2. **AWS CLI installed and configured**
3. **Docker installed**
4. **Basic knowledge of AWS services**

## Step 1: Install and Configure AWS CLI

### Install AWS CLI
```bash
# For Windows (using PowerShell)
winget install -e --id Amazon.AWSCLI

# Or download from: https://aws.amazon.com/cli/
```

### Configure AWS CLI
```bash
aws configure
```
Enter your:
- AWS Access Key ID
- AWS Secret Access Key
- Default region (e.g., us-east-1)
- Default output format (json)

## Step 2: Install Docker Desktop

Download and install Docker Desktop from: https://www.docker.com/products/docker-desktop/

## Step 3: Test Local Build

Before deploying, test your application locally:

```bash
# Build and run with Docker Compose
docker-compose up --build

# Test the API
curl http://localhost:8080/health
```

## Step 4: AWS Infrastructure Setup

### 4.1 Create RDS PostgreSQL Database

1. Go to AWS RDS Console
2. Click "Create database"
3. Choose "Standard create"
4. Select "PostgreSQL"
5. Choose "Free tier" (for testing)
6. Configure:
   - DB instance identifier: `bharat-seva-db`
   - Master username: `postgres`
   - Master password: `YourSecurePassword123!`
   - Public access: Yes (for testing)
   - VPC security group: Create new
7. Click "Create database"

### 4.2 Create ECS Cluster

```bash
aws ecs create-cluster --cluster-name bharat-seva-cluster --region us-east-1
```

### 4.3 Create ECR Repository and Push Image

Run the deployment script:
```bash
# Make script executable (if on Linux/Mac)
chmod +x deploy.sh

# Run deployment script
./deploy.sh
```

## Step 5: Create ECS Task Definition

Create a file called `task-definition.json`:

```json
{
  "family": "bharat-seva-task",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::YOUR_ACCOUNT_ID:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "bharat-seva-app",
      "image": "YOUR_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/bharat-seva-space:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "DB_HOST",
          "value": "YOUR_RDS_ENDPOINT"
        },
        {
          "name": "DB_PORT",
          "value": "5432"
        },
        {
          "name": "DB_USER",
          "value": "postgres"
        },
        {
          "name": "DB_PASSWORD",
          "value": "YourSecurePassword123!"
        },
        {
          "name": "DB_NAME",
          "value": "bharat_seva_space"
        },
        {
          "name": "JWT_SECRET",
          "value": "your-super-secret-jwt-key-change-this-in-production"
        },
        {
          "name": "JWT_EXPIRY",
          "value": "24h"
        },
        {
          "name": "PORT",
          "value": "8080"
        },
        {
          "name": "ENV",
          "value": "production"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/bharat-seva-task",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

Register the task definition:
```bash
aws ecs register-task-definition --cli-input-json file://task-definition.json
```

## Step 6: Create Application Load Balancer

1. Go to EC2 Console → Load Balancers
2. Click "Create load balancer"
3. Choose "Application Load Balancer"
4. Configure:
   - Name: `bharat-seva-alb`
   - Scheme: Internet-facing
   - IP address type: IPv4
   - VPC: Default VPC
   - Mappings: Select all availability zones
5. Configure Security Groups:
   - Allow HTTP (80) and HTTPS (443)
6. Configure Routing:
   - Target group name: `bharat-seva-tg`
   - Target type: IP addresses
   - Port: 8080
   - Protocol: HTTP
7. Click "Create load balancer"

## Step 7: Create ECS Service

```bash
aws ecs create-service \
  --cluster bharat-seva-cluster \
  --service-name bharat-seva-service \
  --task-definition bharat-seva-task \
  --desired-count 1 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxxxx,subnet-yyyyy],securityGroups=[sg-xxxxx],assignPublicIp=ENABLED}" \
  --load-balancers "targetGroupArn=arn:aws:elasticloadbalancing:us-east-1:YOUR_ACCOUNT_ID:targetgroup/bharat-seva-tg,containerName=bharat-seva-app,containerPort=8080"
```

## Step 8: Update Security Groups

1. Go to RDS Console
2. Select your database
3. Go to "Connectivity & security"
4. Click on the security group
5. Add inbound rule:
   - Type: PostgreSQL
   - Port: 5432
   - Source: Custom (ECS security group)

## Step 9: Test Your Deployment

1. Get your ALB DNS name from the EC2 Console
2. Test your API:
```bash
curl http://YOUR_ALB_DNS_NAME/health
```

## Step 10: Set Up Domain (Optional)

1. Register a domain in Route 53 or use existing domain
2. Create hosted zone
3. Create A record pointing to your ALB
4. Configure SSL certificate in Certificate Manager

## Monitoring and Logs

- **ECS Logs**: CloudWatch Logs
- **Application Metrics**: CloudWatch Metrics
- **Database Monitoring**: RDS Console

## Cost Optimization

- Use Spot instances for non-production
- Set up auto-scaling based on CPU/memory
- Use reserved instances for production
- Monitor costs in AWS Cost Explorer

## Security Best Practices

1. Use AWS Secrets Manager for sensitive data
2. Enable VPC Flow Logs
3. Use WAF for additional protection
4. Regular security updates
5. Enable CloudTrail for audit logs

## Troubleshooting

### Common Issues:
1. **Container fails to start**: Check logs in CloudWatch
2. **Database connection issues**: Verify security groups
3. **Load balancer health checks failing**: Check application health endpoint
4. **High costs**: Monitor resource usage and optimize

### Useful Commands:
```bash
# Check ECS service status
aws ecs describe-services --cluster bharat-seva-cluster --services bharat-seva-service

# View logs
aws logs describe-log-groups --log-group-name-prefix /ecs/bharat-seva-task

# Scale service
aws ecs update-service --cluster bharat-seva-cluster --service bharat-seva-service --desired-count 2
```

## Next Steps

1. Set up CI/CD pipeline with GitHub Actions
2. Configure monitoring and alerting
3. Implement auto-scaling
4. Set up backup strategies
5. Configure SSL/TLS certificates

## Support

If you encounter issues:
1. Check AWS CloudWatch logs
2. Review security group configurations
3. Verify environment variables
4. Test locally with Docker Compose first

Good luck with your deployment! 🚀 