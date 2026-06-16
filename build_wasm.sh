#!/bin/bash
# Build Go WASM binary for Breakout game
cd "$(dirname "$0")/go"
echo "Building Go WASM binary..."
env GOOS=js GOARCH=wasm go build -o ../static/wasm/main.wasm
if [ $? -eq 0 ]; then
    echo ""
    echo "Build successful! WASM binary saved to static/wasm/main.wasm"
    echo ""
    echo "Copying wasm_exec.js..."
    GOROOT=$(go env GOROOT)
    if [ -f "$GOROOT/misc/wasm/wasm_exec.js" ]; then
        cp "$GOROOT/misc/wasm/wasm_exec.js" "../static/wasm/wasm_exec.js"
        echo "wasm_exec.js copied successfully."
    else
        echo "Warning: wasm_exec.js not found at $GOROOT/misc/wasm/"
    fi
else
    echo ""
    echo "Build failed! Check error messages above."
    exit 1
fi
