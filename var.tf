variable "build" {}
variable "region" {}

locals {
  instance_type       = "t2.micro"
  app_name            = "hello"
  app_port            = 5000
  app_size            = 5
  app_memory          = 128
  trusted             = "127.0.0.1/32"
  cluster_size        = 2
  redis_instance_type = "cache.t2.micro"
  redis_port          = 6379
  redis_size          = 2
}
