// Generic loader for Go-compiled WASM modules.
// Exposes `loadGoWasm(wasmUrl)` which returns a Promise that resolves when the Go runtime has started.
window.loadGoWasm = function(wasmUrl) {
    if (typeof Go === 'undefined') {
        return Promise.reject(new Error('wasm_exec.js (Go runtime) not loaded. Copy it to static/wasm/wasm_exec.js'))
    }
    const go = new Go();
    // Try instantiateStreaming first for better perf; fall back if server doesn't serve correct mime-type.
    if (WebAssembly.instantiateStreaming) {
        return WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject)
            .then(result => go.run(result.instance))
            .catch(async (err) => {
                // fallback
                const resp = await fetch(wasmUrl);
                const bytes = await resp.arrayBuffer();
                const res = await WebAssembly.instantiate(bytes, go.importObject);
                return go.run(res.instance);
            });
    }
    return fetch(wasmUrl)
        .then(r => r.arrayBuffer())
        .then(bytes => WebAssembly.instantiate(bytes, go.importObject))
        .then(result => go.run(result.instance));
};
