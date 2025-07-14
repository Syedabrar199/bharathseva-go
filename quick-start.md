# 🚀 Quick Start - FREE AWS Deployment

**Total Cost: $0** (AWS Free Tier)

## ⚡ 5-Minute Deployment

### Step 1: Install Tools
```bash
# Install AWS CLI (Windows)
winget install -e --id Amazon.AWSCLI

# Install Docker Desktop
# Download from: https://www.docker.com/products/docker-desktop/
```

### Step 2: Configure AWS
```bash
aws configure
```
Enter your AWS credentials when prompted.

### Step 3: Test Locally (Optional)
```bash
# Test your app works
docker-compose up --build

# Test API
curl http://localhost:8080/health
```

### Step 4: Deploy to AWS (FREE)
```bash
# Make script executable
chmod +x simple-deploy.sh

# Deploy everything
./simple-deploy.sh
```

### Step 5: Test Your Deployment
```bash
# Get your app URL
aws cloudformation describe-stacks \
  --stack-name bharat-seva-free \
  --query "Stacks[0].Outputs[?OutputKey=='LoadBalancerDNS'].OutputValue" \
  --output text

# Test your API
curl http://YOUR_ALB_DNS/health
```

## 🎯 What You Get

- ✅ **Your app running**: `http://your-alb-dns`
- ✅ **Database**: PostgreSQL on RDS
- ✅ **Load Balancer**: For high availability
- ✅ **Monitoring**: CloudWatch logs
- ✅ **Cost**: $0 (Free Tier)

## 🆓 Free Tier Limits

- **EC2 t2.micro**: 750 hours/month (24/7)
- **RDS db.t3.micro**: 750 hours/month
- **Load Balancer**: FREE for 12 months
- **Storage**: 20GB EBS
- **Data Transfer**: 15GB/month

## 🧹 Clean Up (When Done Testing)

```bash
# Delete everything
aws cloudformation delete-stack --stack-name bharat-seva-free

# Verify deletion
aws cloudformation describe-stacks --stack-name bharat-seva-free
```

## 🔧 Troubleshooting

### If deployment fails:
1. Check AWS credentials: `aws sts get-caller-identity`
2. Verify region: `aws configure list`
3. Check CloudFormation events: `aws cloudformation describe-stack-events --stack-name bharat-seva-free`

### If app doesn't respond:
1. Wait 5-10 minutes for full deployment
2. Check EC2 logs: `aws logs describe-log-groups`
3. SSH to EC2: `ssh -i bharat-seva-key.pem ec2-user@YOUR_EC2_IP`

## 📊 Monitor Your App

- **CloudWatch Logs**: Application logs
- **EC2 Console**: Instance status
- **RDS Console**: Database metrics
- **Load Balancer**: Health checks

## 🎉 Success!

Your Go application is now running on AWS for FREE! 

**Next steps:**
1. Test all your API endpoints
2. Monitor performance
3. Set up alerts (optional)
4. Add domain name (optional)

**Remember**: This is for testing. For production, consider:
- SSL certificates
- Auto-scaling
- Backup strategies
- Monitoring alerts 