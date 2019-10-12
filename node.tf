resource "aws_key_pair" "key" {
  key_name   = "node"
  public_key = file(local.public_key)
}

data "aws_ami" "ami" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-ecs-hvm-*-x86_64-ebs"]
  }
}

resource "aws_instance" "node" {
  count                       = 2
  instance_type               = local.instance_type
  ami                         = data.aws_ami.ami.image_id
  security_groups             = [aws_security_group.node.name]
  key_name                    = aws_key_pair.key.id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.ecs-instance-profile.id

  tags = {
    Name = "node"
  }

  connection {
    host = coalesce(self.public_ip, self.private_ip)
    type = "ssh"
    user = "ec2-user"
  }

  provisioner "remote-exec" {
    inline = [
      "echo ECS_CLUSTER=ecs | sudo tee /etc/ecs/ecs.config",
    ]
  }
}
