resource "aws_lb" "alb" {
  name               = "alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = aws_subnet.subnet.*.id
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

resource "aws_lb_target_group" "tg" {
  name     = "tg"
  port     = 5000
  protocol = "HTTP"
  vpc_id   = aws_vpc.vpc.id
  health_check {
    matcher = "200"
    path    = "/hc"
  }
  depends_on = [aws_lb.alb]
}
