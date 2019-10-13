data "aws_ami" "ami" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-*-x86_64-ebs"]
  }
}

resource "aws_key_pair" "key" {
  key_name   = "node"
  public_key = file(local.public_key)
}

resource "aws_launch_configuration" "launch" {
  image_id             = data.aws_ami.ami.image_id
  instance_type        = local.instance_type
  user_data            = <<USERDATA
#!/bin/bash
echo ECS_CLUSTER=ecs | sudo tee /etc/ecs/ecs.config
USERDATA
  key_name             = aws_key_pair.key.id
  security_groups      = [aws_security_group.node.id]
  iam_instance_profile = aws_iam_instance_profile.ecs-instance-profile.id
  # the below is acceptable for the sake of simplicity
  # for more secure production setup vpc could be splitted into private/public subnets
  associate_public_ip_address = true
}

resource "aws_autoscaling_group" "autoscaling" {
  name                 = "asg-${aws_launch_configuration.launch.name}"
  vpc_zone_identifier  = [local.subnet]
  launch_configuration = aws_launch_configuration.launch.name
  min_size             = local.cluster_size
  max_size             = local.cluster_size
  depends_on           = [aws_route.route]
  tag {
    key                 = "Name"
    value               = "node"
    propagate_at_launch = true
  }
}
