
data "archive_file" "lambda_zip" {
  type        = "zip"
  output_path = "./build/lambda/${local.lambda_name}.zip"

  source_file = "${path.root}/bootstrap"
}

resource "aws_s3_object" "lambda_s3_object" {
  bucket      = aws_s3_bucket.lambdas_bucket.bucket
  key         = "${local.lambda_name}.zip"
  source      = data.archive_file.lambda_zip.output_path
  source_hash = data.archive_file.lambda_zip.output_base64sha256
}
