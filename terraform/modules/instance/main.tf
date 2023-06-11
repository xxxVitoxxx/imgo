resource "aws_instance" "main" {
  ami                     = var.ami
  instance_type           = var.instance_type
  subnet_id               = var.subnet_id
  key_name                = var.key_name
  vpc_security_group_ids  = var.vpc_security_group_ids
  disable_api_termination = var.disable_api_termination
  monitoring              = var.monitoring
  user_data               = var.user_data

  root_block_device {
    encrypted   = true
    volume_size = var.volume_size
    volume_type = var.volume_type
  }

  tags = var.tags
}
