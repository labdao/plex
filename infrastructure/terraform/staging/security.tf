resource "aws_security_group" "labdao_public_ssh" {
  name        = "labdao-${var.environment}-sg-ssh-all"
  description = "SSH SG"

  ingress {
    description      = "SSH from anywhere"
    from_port        = 22
    to_port          = 22
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "labdao-${var.environment}-sg-ssh-all"
  }
}

resource "aws_security_group" "labdao_public_bacalhau" {
  name        = "labdao-${var.environment}-sg-bacalhau"
  description = "Public Bacalhau SG"

  ingress {
    description      = "Bacalhau port"
    from_port        = 1234
    to_port          = 1234
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "labdao-${var.environment}-sg-bacalhau"
  }
}

resource "aws_security_group" "labdao_public_plex" {
  name        = "labdao-${var.environment}-sg-plex"
  description = "Public Ports SG"

  ingress {
    description      = "Bacalhau port"
    from_port        = 1234
    to_port          = 1234
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    description      = "IPFS port"
    from_port        = 5001
    to_port          = 5001
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "labdao-${var.environment}-sg-plex"
  }
}

resource "aws_security_group" "labdao_private" {
  name        = "labdao-${var.environment}-sg-private"
  description = "allow all internal traffic"

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [for s in data.aws_subnet.default : s.cidr_block]
    self        = true
  }

  tags = {
    Name = "labdao-${var.environment}-sg-private"
  }
}

resource "aws_security_group" "labdao_egress_all" {
  name        = "labdao-${var.environment}-sg-egress-all"
  description = "Public SG"

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "labdao-${var.environment}-sg-egress-all"
  }
}
