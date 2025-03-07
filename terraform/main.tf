resource "aws_dynamodb_table" "example" {
  name         = local.name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ID"

  attribute {
    name = "ID"
    type = "S"
  }

  global_secondary_index {
    name            = "ModelIndex"
    hash_key        = "model"
    projection_type = "ALL"
  }

  attribute {
    name = "model"
    type = "S"
  }

  stream_enabled   = true
  stream_view_type = "NEW_AND_OLD_IMAGES"

  tags = {
    Name = local.name
  }
}
