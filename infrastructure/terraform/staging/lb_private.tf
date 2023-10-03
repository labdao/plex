# Setting up private LB for compute nodes be able to reach requester and everything to reach IPFS nodes
resource "aws_lb" "labdao_requester_private" {
  name                             = "labdao-requester-${var.environment}"
  internal                         = true
  load_balancer_type               = "network"
  ip_address_type                  = "ipv4"
  subnets                          = data.aws_subnets.default_filtered.ids
  enable_cross_zone_load_balancing = true

}

# Listener for Baclhau API Swarm port for Requester
resource "aws_lb_listener" "labdao_requester_bacalhau_swarm_private_1234" {
  load_balancer_arn = aws_lb.labdao_requester_private.arn
  port              = "1234"
  protocol          = "TCP_UDP"

  # default forward to Bacalhau Swarm TG
  default_action {
    target_group_arn = aws_lb_target_group.labdao_requester_bacalhau_swarm_tg.arn
    type             = "forward"
  }
}

# Listener for Swarm Baclhau Swarm port for Requester
resource "aws_lb_listener" "labdao_requester_bacalhau_swarm_private_1235" {
  load_balancer_arn = aws_lb.labdao_requester_private.arn
  port              = "1235"
  protocol          = "TCP_UDP"

  # default forward to Bacalhau Swarm TG
  default_action {
    target_group_arn = aws_lb_target_group.labdao_requester_bacalhau_swarm_tg.arn
    type             = "forward"
  }
}

# Listener for IPFS Swarm port for IPFS node
resource "aws_lb_listener" "labdao_requester_private_4001" {
  load_balancer_arn = aws_lb.labdao_requester_private.arn
  port              = "4001"
  protocol          = "TCP_UDP"

  # default forward to IPFS Swarm TG
  default_action {
    target_group_arn = aws_lb_target_group.labdao_ipfs_swarm_tg.arn
    type             = "forward"
  }
}

# Listener for IPFS API port for IPFS nodes
resource "aws_lb_listener" "labdao_requester_private_5001" {
  load_balancer_arn = aws_lb.labdao_requester_private.arn
  port              = "5001"
  protocol          = "TCP_UDP"

  # default forward to IPFS API TG
  default_action {
    target_group_arn = aws_lb_target_group.labdao_ipfs_api_tg.arn
    type             = "forward"
  }
}

# TG for Bacalhau API port
resource "aws_lb_target_group" "labdao_requester_bacalhau_api_tg" {
  name     = "labdao-${var.environment}-rqstr-bcl-api"
  port     = 1234
  protocol = "TCP_UDP"

  # PUT in default VPC for now
  vpc_id = data.aws_vpc.default.id

  health_check {
    interval            = 5
    port                = "traffic-port"
    protocol            = "TCP"
    timeout             = 2
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

# TG for Bacalhau Swarm port
resource "aws_lb_target_group" "labdao_requester_bacalhau_swarm_tg" {
  name     = "labdao-${var.environment}-rqstr-bcl-swarm"
  port     = 1235
  protocol = "TCP_UDP"

  # PUT in default VPC for now
  vpc_id = data.aws_vpc.default.id

  health_check {
    interval            = 5
    port                = "traffic-port"
    protocol            = "TCP"
    timeout             = 2
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

# TG for IPFS Swarm port
resource "aws_lb_target_group" "labdao_ipfs_swarm_tg" {
  name     = "labdao-${var.environment}-ipfs-swarm"
  port     = 4001
  protocol = "TCP_UDP"

  # PUT in default VPC for now
  vpc_id = data.aws_vpc.default.id

  health_check {
    interval            = 5
    port                = "traffic-port"
    protocol            = "TCP"
    timeout             = 2
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

# TG for IPFS API port
resource "aws_lb_target_group" "labdao_ipfs_api_tg" {
  name     = "labdao-${var.environment}-ipfs-api"
  port     = 5001
  protocol = "TCP_UDP"

  # PUT in default VPC for now
  vpc_id = data.aws_vpc.default.id

  health_check {
    interval            = 5
    port                = "traffic-port"
    protocol            = "TCP"
    timeout             = 2
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

# private dns record for requester
resource "cloudflare_record" "labdao_requester_private" {
  zone_id = var.cloudflare_zone_id
  name    = "requester.${var.environment}"
  value   = aws_lb.labdao_requester_private.dns_name
  type    = "CNAME"
  ttl     = 3600
}

# private dns record for ipfs
resource "cloudflare_record" "labdao_ipfs_private" {
  zone_id = var.cloudflare_zone_id
  name    = "ipfs.${var.environment}"
  value   = aws_lb.labdao_requester_private.dns_name
  type    = "CNAME"
  ttl     = 3600
}
