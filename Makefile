.PHONY: build
.EXPORT_ALL_VARIABLES:

BUILD_DIR=$(PWD)/terraform/bootstrap
LAMBDA_ROOT=./cmd/ingest
AWS_PROFILE=localstack

infra:
	@echo "Setup Infra"
	docker compose up -d

infra-down:
	@echo "Clean Infra"
	docker compose down

infra-logs:
	@echo "Watch logs"
	docker compose logs -f

tf: build infra
	@echo "Deploying infra"
	cd terraform && terraform init && terraform plan && terraform apply -auto-approve

gen: tf
	@echo "Generating DDB items"
	go run cmd/gen/main.go

build:
	@echo "Building Go Lambda function"
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR) $(LAMBDA_ROOT)/main.go
