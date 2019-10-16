resource "aws_security_group" "node" {
  name        = "node"
  description = "access to node"
  vpc_id      = aws_vpc.vpc.id
  //debug
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [local.trusted]
  }

  ingress {
    from_port       = 32768
    to_port         = 65535
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "alb" {
  name        = "alb"
  description = "access to alb"
  vpc_id      = aws_vpc.vpc.id

  ingress {
    from_port   = local.app_port
    to_port     = local.app_port
    protocol    = "tcp"
    cidr_blocks = [local.trusted]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "redis" {
  name        = "redis"
  description = "access to redis"
  vpc_id      = aws_vpc.vpc.id

  ingress {
    from_port       = local.redis_port
    to_port         = local.redis_port
    protocol        = "TCP"
    security_groups = [aws_security_group.node.id]
  }
}
