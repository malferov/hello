resource "aws_security_group" "node" {
  name        = "node"
  description = "access to node"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [local.trusted]
  }

  ingress {
    from_port   = local.port
    to_port     = local.port
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
