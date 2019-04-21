#!/bin/bash

set -e
set -o pipefail

if [ ! -e tls/ca.pem ]; then
  echo "Generating CA certificate..."
  cfssl gencert -initca tls/ca-csr.json | cfssljson -bare tls/ca
fi

echo "Generating server certificate..."
cfssl gencert -ca=tls/ca.pem -ca-key=tls/ca-key.pem tls/server-csr.json | cfssljson -bare tls/server

echo "Copying certificate files..."
cp tls/server.pem tls/server-key.pem secrets/

echo "Generating CABundle patch..."
CA_BUNDLE=$(cat tls/ca.pem | base64)
sed -i "" -E "s/caBundle: \".*\"/caBundle: \"${CA_BUNDLE}\"/g" patches/webhook-cabundle.yaml
