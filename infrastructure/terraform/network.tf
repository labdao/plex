resource "aws_security_group" "plex" {
  name = "dev-web"
  description = "SSH with key and HTTP open"
}

resource "aws_security_group" "internal" {
  name = "receptor-web"
  description = "allow all internal traffic"
  ingress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    self = true
  }
  // allow all egress
  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}
