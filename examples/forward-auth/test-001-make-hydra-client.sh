#!/bin/sh

set -e

echo ""
echo "Deleting old client if it exists"
docker-compose exec hydra \
  hydra clients delete \
  --endpoint http://127.0.0.1:4445 \
  auth-code-client || true

echo ""
echo "Creating new client"
docker-compose exec hydra \
  hydra clients create \
  --endpoint http://127.0.0.1:4445 \
  --id auth-code-client \
  --secret secret \
  --grant-types authorization_code,refresh_token \
  --response-types code,id_token \
  --scope openid,offline,profile,email \
  --callbacks http://auth-meso.test.localhost/_oauth

echo ""
echo "New client created!"
echo "Client ID: auth-code-client"
echo "Secret:    secret"
