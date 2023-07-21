# Setting up private LB for compute nodes be able to reach requester
resource "aws_lb" "labdao_requester_private" {
  name               = "labdao-requester-${var.environment}"
  internal           = true
  load_balancer_type = "network"
  ip_address_type    = "ipv4"
  subnets            = data.aws_subnets.default_filtered.ids
}

resource "aws_lb_listener" "labdao_requester_bacalhau_swarm_private_1235" {
  load_balancer_arn = aws_lb.labdao_requester_private.arn
  port              = "1235"
  protocol          = "TCP_UDP"

  default_action {
    target_group_arn = aws_lb_target_group.labdao_requester_bacalhau_swarm_tg.arn
    type             = "forward"
  }
}

resource "aws_lb_listener" "labdao_requester_private_4001" {
  load_balancer_arn = aws_lb.labdao_requester_private.arn
  port              = "4001"
  protocol          = "TCP_UDP"

  default_action {
    target_group_arn = aws_lb_target_group.labdao_requester_ipfs_swarm_tg.arn
    type             = "forward"
  }
}

resource "aws_lb_target_group" "labdao_requester_bacalhau_swarm_tg" {
  name     = "labdao-${var.environment}-rqstr-bcl-swarm"
  port     = 1235
  protocol = "TCP_UDP"

  # PUT in default VPC for now
  vpc_id = data.aws_vpc.default.id

  # deregistration_delay = 5

  # slow_start = 60

  health_check {
    interval            = 5
    port                = "traffic-port"
    protocol            = "TCP"
    timeout             = 2
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

resource "aws_lb_target_group" "labdao_requester_ipfs_swarm_tg" {
  name     = "labdao-${var.environment}-rqstr-ipfs-swarm"
  port     = 4001
  protocol = "TCP_UDP"

  # PUT in default VPC for now
  vpc_id = data.aws_vpc.default.id

  # deregistration_delay = 5

  # slow_start = 60

  health_check {
    interval            = 5
    port                = "traffic-port"
    protocol            = "TCP"
    timeout             = 2
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

# public dns record for recepter
resource "cloudflare_record" "labdao_requester_private" {
  zone_id = var.cloudflare_zone_id
  name    = "requester.${var.environment}"
  value   = aws_lb.labdao_requester_private.dns_name
  type    = "CNAME"
  ttl     = 3600
}
