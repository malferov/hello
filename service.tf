resource "aws_ecs_cluster" "ecs" {
  name = "ecs"
}

resource "aws_ecs_service" "srv" {
  name            = local.app
  cluster         = aws_ecs_cluster.ecs.id
  task_definition = aws_ecs_task_definition.task.arn
  desired_count   = 1
  load_balancer {
    target_group_arn = aws_lb_target_group.tg.arn
    container_name   = local.app
    container_port   = local.port
  }
}

resource "aws_ecs_task_definition" "task" {
  family                = "service"
  container_definitions = <<TASK
[
  {
    "name": "hello",
    "image": "docker.io/malferov/app:${var.build}",
    "memory": 128,
    "essential": true,
    "portMappings": [
      {
        "containerPort": ${local.port},
        "hostPort": 0
      }
    ]
  }
]
TASK
}
