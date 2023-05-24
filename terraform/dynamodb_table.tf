terraform {
  required_version = ">= 1.0.0"
}

resource "aws_dynamodb_table" "main" {
  name         = "terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}
