output "vpc_id" {
  value       = aws_vpc.main.id
  description = "aws devops vpc id"
}

output "bucket_name" {
  value       = aws_s3_bucket.main.id
  description = "aws devops bucket domain"
}

output "dynamo_id" {
  value       = aws_dynamodb_table.main.id
  description = "aws devops dynamo table name"
}

output "load_balalnce_security_group_id" {
    value = aws_security_group.load_balancer_security_group_external.id
}

output "instance_security_group_id" {
    value = aws_security_group.instance_security_group_internal.id
}

output "load_balancer_id" {
    value = aws_lb.main.id
}
