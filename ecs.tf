resource "aws_ecs_cluster" "ecs" {
  name = "ecs"
}

resource "aws_ecs_service" "srv" {
  name            = "hello"
  cluster         = aws_ecs_cluster.ecs.id
  task_definition = aws_ecs_task_definition.task.arn
  desired_count   = 1
}

resource "aws_ecs_task_definition" "task" {
  family                = "service"
  container_definitions = file("task-definition.json")
}
