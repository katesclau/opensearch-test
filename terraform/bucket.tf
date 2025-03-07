resource "aws_s3_bucket" "lambdas_bucket" {
  bucket = "${local.name}-lambda"
}

resource "aws_s3_bucket_public_access_block" "lambdas_bucket_public_access_block" {
  bucket = aws_s3_bucket.lambdas_bucket.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_versioning" "lambdas_bucket_versioning" {
  bucket = aws_s3_bucket.lambdas_bucket.id

  versioning_configuration {
    status = "Enabled"
  }
}
