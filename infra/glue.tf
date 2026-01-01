resource "aws_glue_catalog_database" "logs" {
  name = "logs"
}

resource "aws_glue_catalog_table" "aws_glue_catalog" {
  name          = "table_name"
  database_name = aws_glue_catalog_database.logs.name

  storage_descriptor {
    location      = "s3://rvsn-logs/"
    input_format  = "org.apache.hadoop.mapred.TextInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat"

    ser_de_info {
      name                  = "json"
      serialization_library = "org.openx.data.jsonserde.JsonSerDe"
    }

    columns {
      name    = "id"
      type    = "string"
      comment = "the user id"
    }
  }
}
