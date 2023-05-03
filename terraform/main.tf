terraform {
  required_version = "~> 1.0"
}

# Ensure default profile is configured or AWS_DEFAULT_PROFILE is set
provider "aws" {
  region = "ap-southeast-1"
}

locals {
  tags = {
    Service     = "craftctl"
    Environment = "test"
  }
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "aws_s3_bucket" "b" {
  bucket_prefix = "test-"
  force_destroy = true
  tags          = local.tags
}

resource "aws_iam_user" "u" {
  name = "test-${random_string.suffix.result}"
  path = "/${local.tags.Service}/"
  tags = local.tags
}

resource "aws_iam_user_policy" "up" {
  name = "s3-access"
  user = aws_iam_user.u.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action   = ["s3:PutObject"]
        Effect   = "Allow"
        Resource = ["${aws_s3_bucket.b.arn}/*"]
      },
    ]
  })
}

resource "aws_iam_access_key" "k" {
  user = aws_iam_user.u.name
}

output "bucket" {
  value = aws_s3_bucket.b.bucket
}

output "aws_access_key_id" {
  value = aws_iam_access_key.k.id
}

output "aws_secret_access_key" {
  sensitive = true
  value     = aws_iam_access_key.k.secret
}
