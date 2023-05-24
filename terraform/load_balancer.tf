terraform {
  required_version = ">= 1.0.0"
}

resource "aws_lb" "main" {
  name               = "load-balancer-external"
  load_balancer_type = "application"
  internal           = false
  security_groups    = [aws_security_group.load_balancer_security_group_external.id]
  subnets = [
    aws_subnet.public_subnet_1a.id,
    aws_subnet.public_subnet_1c.id
  ]
}

resource "aws_lb_target_group" "main" {
  name        = "target-group-instance"
  target_type = "ip"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
}

resource "aws_lb_listener_rule" "main" {
  listener_arn = aws_lb_listener.listen_http_443.arn

  action {
    type = "fixed-response"

    fixed_response {
      status_code  = "200"
      content_type = "text/plain"
      message_body = "hello"
    }
  }

  condition {
    host_header {
      values = ["imgo.vitooooo.com"]
    }
  }
}

data "aws_acm_certificate" "main" {
  domain    = "*.vitooooo.com"
  types     = ["AMAZON_ISSUED"]
  key_types = ["RSA_2048"]
  statuses  = ["ISSUED"]
}

resource "aws_lb_listener" "listen_http_443" {
  load_balancer_arn = aws_lb.main.arn
  port              = 443
  protocol          = "HTTPS"
  certificate_arn   = data.aws_acm_certificate.main.arn

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "text/plain"
      status_code  = 200
      message_body = "page not found"
    }
  }
}
