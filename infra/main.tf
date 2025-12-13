resource "aws_s3_bucket" "rvsnlogs" {
  bucket        = "rvsnlogs"
  force_destroy = true
  tags = {
    Name        = "My logs bucket"
    Environment = "Dev"
  }
}
