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
  default_tags {
    tags = {
      Project = "imgo"
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

resource "aws_subnet" "private_subnet" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.11.2.0/24"
  availability_zone = "ap-northeast-1a"
}

resource "aws_subnet" "public_subnet_1c" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.11.3.0/24"
  availability_zone       = "ap-northeast-1c"
  map_public_ip_on_launch = true
}

data "aws_ami" "main" {
  most_recent = true

  filter {
    name   = "name"
    values = ["imgo*"]
  }
}

module "ec2" {
  source                  = "./modules/instance"
  ami                     = data.aws_ami.main.id
  instance_type           = "t2.micro"
  subnet_id               = aws_subnet.private_subnet.id
  key_name                = "vito"
  vpc_security_group_ids  = [aws_security_group.instance_security_group_internal.id]
  disable_api_termination = true
  volume_size             = 8
  volume_type             = "gp2"
  monitoring              = true
  user_data               = <<-EOF
  #!/bin/bash

  cd /home/ubuntu/build
  sudo nohup ./imgo >/dev/null 2>&1 &
  EOF

  tags = {
    Name : "imgo-instance"
  }
}
