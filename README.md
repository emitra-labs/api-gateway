# api-gateway

A reverse proxy server that acts as a gateway for multiple backend services.

## How It Works

This program reads all environment variables and filters those with the suffix `_HTTP_ADDRESS`. These keys represent backend services to be registered. The values must be valid URLs without trailing slashes. Here are some examples of supported key-value pairs:

```
FOO_HTTP_ADDRESS=http://localhost:8080
BAR_HTTP_ADDRESS=https://my-secure-backend:3000
BAZ_HTTP_ADDRESS=http://baz-service:3000/api
MY_BACKEND_HTTP_ADDRESS=https://my-backend.example.com
```

These services will be exposed based on their service paths. The service path is derived by trimming the `_HTTP_ADDRESS` suffix and converting the remaining part to kebab-case. For example, the service path for `BILLING_PAYMENT_HTTP_ADDRESS` is `/billing-payment`. You can then make a request to `http://my-api-gateway/billing-payment/invoices`.

## Quickstart

You have two options to get started: using the Docker image or running the program locally.

### 1. Using the Docker Image

```bash
docker run \
  --name api-gateway \
  -p 3000:3000 \
  -e FOO_HTTP_ADDRESS=http://foo-service:8080 \
  ghcr.io/ukasyah-dev/api-gateway:main
```

### 2. Running the Program Locally

```bash
# Create .env file, update its content
cp .env.example .env

# Run the app
make run
```
