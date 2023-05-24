terraform {
  required_version = ">= 1.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }

  backend "s3" {
    bucket         = "imgo-bucket"
    key            = "terraform/terraform.tfstate"
    region         = "ap-northeast-1"
    dynamodb_table = "terraform-locks"
    encrypot       = true
  }
}


provider "aws" {
  access_key = var.access_key
  secret_key = var.secret_access_key
  region     = var.region

  default_tags {
    tags = {
      project = "imgo"
    }
  }
}

resource "aws_vpc" "main" {
  cidr_block           = "10.11.0.0/16"
  enable_dns_hostnames = true
}

resource "aws_subnet" "public_subnet_1a" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.11.1.0/24"
  availability_zone       = "ap-northeast-1a"
  map_public_ip_on_launch = true
}

resource "aws_subnet" "public_subnet_1c" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.11.3.0/24"
  availability_zone       = "ap-northeast-1c"
  map_public_ip_on_launch = true
}

resource "aws_subnet" "private_subnet" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.11.2.0/24"
  availability_zone = "ap-northeast-1a"
}

resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route_table" "private_route_table" {
  vpc_id = aws_vpc.main.id
}

resource "aws_route" "public_route" {
  route_table_id         = aws_route_table.public_route_table.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.gateway.id
}

resource "aws_route_table_association" "public_association_1a" {
  subnet_id      = aws_subnet.public_subnet_1a.id
  route_table_id = aws_route_table.public_route_table.id
}

resource "aws_route_table_association" "public_association_1c" {
  subnet_id      = aws_subnet.public_subnet_1c.id
  route_table_id = aws_route_table.public_route_table.id
}

resource "aws_route_table_association" "private_association" {
  subnet_id      = aws_subnet.private_subnet.id
  route_table_id = aws_route_table.private_route_table.id
}


resource "aws_security_group" "load_balancer_security_group_external" {
  name        = "load-balancer-security-group-external"
  description = "allow 443 inbound traffic for external"
  vpc_id      = aws_vpc.main.id

  ingress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    self        = false
    description = "inbound traffic from external with 443 port"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "instance_security_group_internal" {
  name        = "instance-security-group-internal"
  description = "allow 80 inbound traffic for load balancer"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    self            = false
    description     = "inbound traffic from load balancer with 80 port"
    security_groups = [aws_security_group.load_balancer_security_group_external.id]
  }

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    self            = false
    description     = "inbound traffic from load balancer health check with 8080 port"
    security_groups = [aws_security_group.load_balancer_security_group_external.id]
  }
}
