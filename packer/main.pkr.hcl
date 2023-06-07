packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

locals { date = formatdate("YYMMDDhhmmss", timestamp()) }

source "amazon-ebs" "ubuntu" {
  ami_name      = "imgo-${local.date}"
  instance_type = "t2.micro"
  region        = "ap-northeast-1"
  vpc_id        = "vpc-0cf65d471f5dd2ca9"
  subnet_id     = "subnet-0f86d4348004e1bf6"

  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }

    most_recent = true
    owners      = ["099720109477"]
  }
  
  ssh_username = "ubuntu"
}

build {
  name = "imgo-${local.date}"
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell-local" {
    command = "sh build.sh"
  }

  provisioner "file" {
    source      = "../build"
    destination = "."
  }

  provisioner "file" {
    source      = "../config.toml"
    destination = "./build/config.toml"
  }
}
