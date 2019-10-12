/*
resource "aws_lb" "internal_microservices" {
  name                       = "Microservices"
  internal                   = true
  load_balancer_type         = "application"
  security_groups            = [aws_security_group.alb_microservices.id]
  subnets                    = module.vpc.private_subnets
  tags                       = local.common_tags
}
*/
data "aws_vpcs" "vpc" {}

locals {
  vpc = tolist(data.aws_vpcs.vpc.ids)[0]
}

data "aws_subnet_ids" "subnet" {
  vpc_id = local.vpc
}

resource "aws_lb" "alb" {
  name               = "alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = data.aws_subnet_ids.subnet.ids
}

resource "aws_lb_listener" "lsr" {
  load_balancer_arn = aws_lb.alb.arn
  port              = "5000"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.tg.arn
  }
}
/*
resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.external_kong.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:acm:eu-west-1:611887617466:certificate/c31446fe-e867-456a-9dd7-836da2987f5c"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.https_kong.arn
  }
}
*/
/*
resource "aws_lb_target_group" "http" {
  name     = "tg-http"
  port     = 80
  protocol = "HTTP"
  vpc_id   = module.vpc.vpc_id
}
*/
resource "aws_lb_target_group" "tg" {
  name     = "tg"
  port     = 5000
  protocol = "HTTP"
  vpc_id   = local.vpc
  health_check {
    matcher = "200"
    path    = "/"
    #    port     = "5000"
    #    protocol = "HTTP"
  }
}
/*
# Target Group
resource "aws_lb_target_group" "target_group" {
  name = "tg-${var.service_name}"
  tags{
    Name = "${var.service_name}-service"
  }
  port = "${var.microservice_port}"
  protocol = "${var.lb_protocol}"
  vpc_id = "${var.microservices_vpc}"
  health_check {
    path = "${var.healthcheck_path}"
    matcher = "${var.healthcheck_port_matcher}"
  }
  lifecycle {
    create_before_destroy = true
  }
}
*/


/*
resource "aws_lb_target_group_attachment" "tg_kong" {
  target_group_arn = aws_lb_target_group.https_kong.arn
  target_id        = aws_instance.kong.id
  port             = 8443
}

*/