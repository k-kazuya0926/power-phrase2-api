# VPC
resource "aws_vpc" "tf_vpc" {
  cidr_block           = "172.31.0.0/16"
  instance_tenancy     = "default"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    Name = "tf_vpc"
  }
}

# Subnet
resource "aws_subnet" "public_a" {
  vpc_id                  = aws_vpc.tf_vpc.id
  cidr_block              = "172.31.0.0/20"
  availability_zone       = "ap-northeast-1a"
  map_public_ip_on_launch = true
  tags = {
    Name = "tf_public_a"
  }
}

resource "aws_subnet" "public_d" {
  vpc_id                  = aws_vpc.tf_vpc.id
  cidr_block              = "172.31.16.0/20"
  availability_zone       = "ap-northeast-1d"
  map_public_ip_on_launch = true
  tags = {
    Name = "tf_public_d"
  }
}

resource "aws_subnet" "private_a" {
  vpc_id            = aws_vpc.tf_vpc.id
  cidr_block        = "172.31.32.0/20"
  availability_zone = "ap-northeast-1a"
  tags = {
    Name = "tf_private_a"
  }
}

resource "aws_subnet" "private_d" {
  vpc_id            = aws_vpc.tf_vpc.id
  cidr_block        = "172.31.48.0/20"
  availability_zone = "ap-northeast-1d"
  tags = {
    Name = "tf_private_d"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.tf_vpc.id
  tags = {
    Name = "tf-gw"
  }
}

# Route Table
resource "aws_route_table" "public_rtb" {
  vpc_id = aws_vpc.tf_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }
  tags = {
    Name = "tf_rtb"
  }
}

resource "aws_route_table_association" "public_a" {
  subnet_id      = aws_subnet.public_a.id
  route_table_id = aws_route_table.public_rtb.id
}

resource "aws_route_table_association" "public_d" {
  subnet_id      = aws_subnet.public_d.id
  route_table_id = aws_route_table.public_rtb.id
}

# Security Group
resource "aws_security_group" "app" {
  name        = "tf_web"
  description = "It is a security group on http of tf_vpc"
  vpc_id      = aws_vpc.tf_vpc.id
  tags = {
    Name = "tf_web"
  }
}

# Security Group Rule
resource "aws_security_group_rule" "ssh" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.app.id
}

resource "aws_security_group_rule" "all" {
  type              = "egress"
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.app.id
}

resource "aws_security_group" "db" {
  name        = "tf_db"
  description = "It is a security group on db of tf_vpc."
  vpc_id      = aws_vpc.tf_vpc.id
  tags = {
    Name = "tf_db"
  }
}

resource "aws_security_group_rule" "db" {
  type                     = "ingress"
  from_port                = 3306
  to_port                  = 3306
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.app.id
  security_group_id        = aws_security_group.db.id
}

# DB Subnet Group
resource "aws_db_subnet_group" "main" {
  name        = "tf_dbsubnet"
  description = "It is a DB subnet group on tf_vpc."
  subnet_ids  = [aws_subnet.private_a.id, aws_subnet.private_d.id]
  tags = {
    Name = "tf_dbsubnet"
  }
}

# RDS
resource "aws_db_instance" "db" {
  identifier              = "tf-dbinstance"
  allocated_storage       = 10
  engine                  = "mysql"
  engine_version          = "8.0.21"
  instance_class          = "db.t2.micro"
  storage_type            = "gp2"
  name                    = "power_phrase2"
  username                = var.aws_db_username
  password                = var.aws_db_password
  backup_retention_period = 1
  vpc_security_group_ids  = [aws_security_group.db.id]
  db_subnet_group_name    = aws_db_subnet_group.main.name
  skip_final_snapshot     = true
}

# S3
resource "aws_s3_bucket" "power-phrase" {
  acl = "public-read"
  website {
    error_document = "error.html"
    index_document = "index.html"
  }
}

resource "aws_s3_bucket" "myfiles" {
  acl = "private"
}
