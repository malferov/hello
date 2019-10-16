resource "aws_ecs_cluster" "ecs" {
  name = "ecs"
}

resource "aws_ecs_service" "srv" {
  name            = local.app_name
  cluster         = aws_ecs_cluster.ecs.id
  task_definition = aws_ecs_task_definition.task.arn
  desired_count   = local.app_size
  load_balancer {
    target_group_arn = aws_lb_target_group.tg.arn
    container_name   = local.app_name
    container_port   = local.app_port
  }
}

resource "aws_ecs_task_definition" "task" {
  family                = "service"
  container_definitions = <<TASK
[
  {
    "name": "${local.app_name}",
    "image": "${aws_ecr_repository.ecr.repository_url}:${var.build}",
    "memory": 128,
    "essential": true,
    "environment": [
      {
        "name": "GIN_MODE",
	"value": "release"
      },
      {
        "name": "REDIS_ENDPOINT",
	"value": "${aws_elasticache_replication_group.redis.primary_endpoint_address}:${local.redis_port}"
      }
    ],
    "portMappings": [
      {
        "containerPort": ${local.app_port},
        "hostPort": 0
      }
    ]
  }
]
TASK
}
