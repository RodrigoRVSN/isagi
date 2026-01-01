resource "aws_s3_bucket" "rvsnlogs" {
  bucket        = "rvsn-logs"
  force_destroy = true
  tags = {
    Name        = "My logs bucket"
    Environment = "Dev"
  }
}
