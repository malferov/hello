#variable "vpc" {}
#variable "microservices_sg" {}
#variable "instance_type" {}
/*
data "aws_security_group" "microservices" {
  tags = {
    Type = "${var.microservices_sg}"
  }
}

data "aws_vpc" "vpc" {
  tags = {
    Type = "Microservices"
  }
}

data "aws_subnet_ids" "subnet" {
  vpc_id = "${data.aws_vpc.vpc.id}"
  tags = {
    Network = "Private"
  }
}
*/
resource "aws_elasticache_subnet_group" "redis" {
  name        = "redis"
  subnet_ids  = aws_subnet.subnet.*.id
  description = "redis subnet group"
}

resource "aws_elasticache_replication_group" "redis" {
  automatic_failover_enabled    = true
  replication_group_id          = "redis-micro"
  replication_group_description = "redis replication group"
  node_type                     = local.redis_instance_type
  number_cache_clusters         = local.redis_size
  parameter_group_name          = "default.redis5.0"
  engine_version                = "5.0.3"
  security_group_ids            = [aws_security_group.redis.id]
  port                          = local.redis_port
  subnet_group_name             = aws_elasticache_subnet_group.redis.name
}
