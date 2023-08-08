resource "aws_security_group" "plex" {
  name        = "dev-web"
  description = "SSH with key and HTTP open"
}

resource "aws_security_group" "internal" {
  name        = "internal-sg"
  description = "allow all internal traffic"
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    self      = true
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

resource "aws_security_group" "external_ssh" {
  name        = "external-ssh-sg"
  description = "Allow ssh from outside"
}

resource "aws_security_group_rule" "ingress_ssh_external" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  security_group_id = aws_security_group.external_ssh.id
}

resource "aws_security_group" "allow_metabase_postgres" {
  name        = "allow-metabase"
  description = "Allow Metabase IPs to access RDS"
}

resource "aws_security_group_rule" "ingress_allow_metabase_postgres" {
  type              = "ingress"
  from_port         = 5432
  to_port           = 5432
  protocol          = "tcp"
  # https://www.metabase.com/docs/latest/cloud/ip-addresses-to-whitelist
  cidr_blocks       = ["18.207.81.126/32", "3.211.20.157/32", "50.17.234.169/32"]
  security_group_id = aws_security_group.allow_metabase_postgres.id
}
