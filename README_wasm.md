Building the Go -> WASM Breakout demo

Prerequisites
- Go installed (1.11+ with wasm support). Ensure `go env` works.

Steps

1. Copy `wasm_exec.js` from your Go distribution into `static/wasm/`.
   - On Windows the file is usually at `%GOROOT%\\misc\\wasm\\wasm_exec.js`.

2. Build the WASM binary from the `go` directory:

Windows (PowerShell):

```powershell
$env:GOOS = 'js'
$env:GOARCH = 'wasm'
go build -o ../static/wasm/main.wasm
```

Linux/macOS:

```bash
env GOOS=js GOARCH=wasm go build -o ../static/wasm/main.wasm
```

3. Serve the Flask app (or a static server). The page `templates/fun.html` will load `static/wasm/wasm_exec.js`, `static/wasm/bootstrap.js`, and `static/wasm/main.wasm`.

Notes
- Go's WASM runtime is sizeable. Use `go build -ldflags "-s -w"` to strip symbols if you need smaller binaries.
- If `instantiateStreaming` fails due to mime-type issues, `bootstrap.js` falls back to fetching the bytes.
