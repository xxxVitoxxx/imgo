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
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    self            = false
    description     = "inbound traffic from border with 22 port"
    security_groups = [var.border_security_group_id]
  }

  ingress {
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    self            = false
    description     = "inbound traffic from border with 80 port"
    security_groups = [aws_security_group.load_balancer_security_group_external.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
