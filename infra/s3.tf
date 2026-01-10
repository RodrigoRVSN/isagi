resource "aws_s3_bucket" "rvsnlogs" {
  bucket        = "rvsn-logs"
  force_destroy = true
  tags = {
    Name        = "My logs bucket"
    Environment = "Dev"
  }
}

resource "aws_s3_bucket_public_access_block" "logs_access_block" {
  bucket = aws_s3_bucket.rvsnlogs.id

  block_public_acls   = true
  block_public_policy = true
  ignore_public_acls  = true
}
