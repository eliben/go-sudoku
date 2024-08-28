// Web worker for invoking generateBoard from Go in a separate thread.
importScripts("wasm_exec.js");
console.log("Worker is running");

// Load the WASM module with Go code.
const go = new Go();
WebAssembly.instantiateStreaming(fetch("gosudoku.wasm"), go.importObject).then(
    (result) => {
        go.run(result.instance);
        console.log("Worker loaded WASM module");
    }).catch((err) => {
        console.error("Worker failed to load WASM module: ", err)
    });

// The worker's logic is very simple: it waits for a "generate" message with
// paramterers, runs the Go code to generate the board and sends a "ready"
// message with the SVG text back to the main thread.
onmessage = ({ data }) => {
    let { action, payload } = data;
    console.log("Worker received message: ", action, payload);
    switch (action) {
        case "generate":
            let svgText = generateBoard(payload.hint, payload.symmetrical);
            postMessage({ action: "boardReady", payload: svgText });
            break;
        default:
            throw (`unknown action '${action}'`);
    }
};
