# Test OpenSearch and Ingestion Pipeline Setup

## Overview

This project is designed to set up a test environment for an OpenSearch domain and an ingestion pipeline. The primary objective is to fetch data from a DynamoDB (DDB) table and ingest it into the OpenSearch service for indexing and searching.

## Components

### OpenSearch Domain

The OpenSearch domain is configured with the necessary security, logging, and access policies to ensure secure and efficient data handling. It includes settings for instance types, counts, and advanced security options.

### Ingestion Pipeline

The ingestion pipeline is responsible for streaming data from the DynamoDB table to the OpenSearch domain. It uses AWS services like AWS OSIS (OpenSearch Ingestion Service) to facilitate this process.

## Features

- **Secure Data Transfer**: The setup ensures secure data transfer between DynamoDB and OpenSearch using VPC options and security groups.
- **Logging and Monitoring**: CloudWatch log groups are configured to monitor the ingestion process and OpenSearch domain activities.
- **Scalability**: The configuration includes options for autoscaling to handle varying data loads efficiently.
- **Custom Endpoint**: A custom endpoint is enabled for the OpenSearch domain to allow access from the internet with appropriate security measures.

## Usage

This setup can be used for testing and validating the data ingestion process from DynamoDB to OpenSearch. It provides a robust framework to ensure data integrity and security during the ingestion process.
