#!/bin/bash
set -e
apt -qq update && apt -qq install brotli
echo "==> Bundling wrapper"
npm config set update-notifier false
npm install
npx --no-update-notifier webpack
echo "==> Compiling to wasm"
javy -o /tmp/mjml.wasm /tmp/mjml.js
echo "==> Compressing wasm"
brotli -f -o ../wasm/mjml.wasm.br /tmp/mjml.wasm
echo "==> Pkging test server"
npx pkg -o ../node-test-server/server --compress brotli --targets node18-linux-x64 /tmp/server.js
echo "==> Done!"