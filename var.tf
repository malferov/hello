variable "build" {}
variable "region" {}

locals {
  instance_type = "t2.micro"
  public_key    = ".key/id_rsa.pub"
  port          = 5000
  trusted       = "85.148.177.204/32"
  app           = "hello"
  cluster_size  = 1
}
