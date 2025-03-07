resource "aws_iam_role" "lambda_role" {
  name = "${local.name}-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name = "${local.name}-lambda-role"
  }
}

resource "aws_iam_role_policy" "lambda_policy" {
  name = "${local.name}-lambda-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:DescribeStream",
          "dynamodb:GetRecords",
          "dynamodb:GetShardIterator",
          "dynamodb:ListStreams"
        ]
        Resource = aws_dynamodb_table.example.stream_arn
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      }
    ]
  })
}

resource "aws_lambda_function" "ddb_stream_processor" {
  depends_on = [
    aws_s3_object.lambda_s3_object
  ]
  function_name     = local.lambda_name
  timeout           = 60
  memory_size       = 512
  s3_bucket         = aws_s3_bucket.lambdas_bucket.bucket
  s3_key            = aws_s3_object.lambda_s3_object.key
  s3_object_version = aws_s3_object.lambda_s3_object.version_id
  handler           = "index.handler"
  runtime           = "provided.al2023"
  role              = aws_iam_role.lambda_role.arn

  environment {
    variables = {
      OPENSEARCH_URL = "https://opensearch-node1:9200"
    }
  }

  tags = {
    Name = local.lambda_name
  }
}

resource "aws_lambda_event_source_mapping" "ddb_stream" {
  event_source_arn  = aws_dynamodb_table.example.stream_arn
  function_name     = aws_lambda_function.ddb_stream_processor.arn
  starting_position = "LATEST"
}
