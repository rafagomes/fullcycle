#!/bin/bash

echo "Teste de Rate Limiter via IP (sem token)"
for i in {1..6}; do
  response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/)
  echo "Requisição $i: $response"
done

echo ""
echo "Teste de Rate Limiter via Token"
for i in {1..11}; do
  response=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: abc123" http://localhost:8080/)
  echo "Requisição $i: $response"
done