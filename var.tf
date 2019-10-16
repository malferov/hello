variable "build" {}
variable "region" {}

locals {
  instance_type       = "t2.micro"
  app_name            = "hello"
  app_port            = 5000
  app_size            = 1
  trusted             = "85.148.177.204/32"
  cluster_size        = 1
  redis_instance_type = "cache.t2.micro"
  redis_port          = 6379
  redis_size          = 2
}
