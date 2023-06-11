output "vpc_id" {
  description = "ID of the specific VPC to retrieve"
  value       = aws_vpc.main.id
}

output "bucket_name" {
  description = "Name of the bucket"
  value       = aws_s3_bucket.main.id
}

output "dynamo_name" {
  description = "Name of the DynamoDB table"
  value       = aws_dynamodb_table.main.id
}

output "load_balalncer_security_group_id" {
  description = "ID of the specific security group of Load Balancer to retrieve"
  value       = aws_security_group.load_balancer_security_group_external.id
}

output "instance_security_group_id" {
  description = "ID of the specific security group of EC2 instance to retrieve"
  value       = aws_security_group.instance_security_group_internal.id
}
