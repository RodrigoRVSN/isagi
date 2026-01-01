resource "aws_s3_bucket" "athena_results" {
  bucket = "rvsn-athena-results"
}

resource "aws_athena_database" "athena_db" {
  name   = "database_name"
  bucket = aws_s3_bucket.athena_results.id
}

resource "aws_athena_workgroup" "workgroup" {
  name = "workgroup_name"

  configuration {
    result_configuration {
      output_location = "s3://${aws_s3_bucket.athena_results.bucket}"
    }
  }
}

resource "aws_athena_prepared_statement" "test" {
  name            = "statement_test"
  query_statement = "SELECT * FROM ${aws_athena_database.athena_db.name}"
  workgroup       = aws_athena_workgroup.workgroup.name
}
