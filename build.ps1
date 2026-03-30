$ErrorActionPreference = "Stop"

Write-Host "Building Ninesweeper for WASM..." -ForegroundColor Cyan

# Create build directory
New-Item -ItemType Directory -Force -Path build | Out-Null

# Build WASM binary
$env:GOOS = "js"
$env:GOARCH = "wasm"
go build -o build/ninesweeper.wasm .
Remove-Item Env:GOOS
Remove-Item Env:GOARCH

# Copy wasm_exec.js from Go installation
$goRoot = (go env GOROOT)
$wasmExecPaths = @(
    "$goRoot\misc\wasm\wasm_exec.js",
    "$goRoot\lib\wasm\wasm_exec.js"
)
$found = $false
foreach ($p in $wasmExecPaths) {
    if (Test-Path $p) {
        Copy-Item $p build\
        Write-Host "Found wasm_exec.js at $p" -ForegroundColor Green
        $found = $true
        break
    }
}
if (-not $found) {
    Write-Host "ERROR: Could not find wasm_exec.js" -ForegroundColor Red
    Write-Host "Searched:" -ForegroundColor Red
    foreach ($p in $wasmExecPaths) { Write-Host "  $p" }
    Write-Host ""
    Write-Host "Run this to find it manually:" -ForegroundColor Yellow
    Write-Host "  Get-ChildItem -Path (go env GOROOT) -Recurse -Filter wasm_exec.js"
    exit 1
}

# Generate index.html
@"
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Ninesweeper</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            background: #1a1a2e;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
            font-family: 'Segoe UI', sans-serif;
            color: #eee;
        }
        h1 {
            margin-bottom: 8px;
            font-size: 1.4em;
            letter-spacing: 2px;
        }
        p {
            margin-bottom: 16px;
            font-size: 0.85em;
            color: #aaa;
        }
        canvas {
            border: 2px solid #333;
            border-radius: 4px;
        }
        #loading {
            font-size: 1.1em;
            color: #FFD700;
            margin-top: 20px;
        }
    </style>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(
            fetch('ninesweeper.wasm'),
            go.importObject
        ).then(result => {
            document.getElementById('loading').style.display = 'none';
            go.run(result.instance);
        });
    </script>
</head>
<body>
    <h1>NINESWEEPER</h1>
    <p>Sudoku meets Minesweeper &mdash; right click to flag, numbers to fill</p>
    <div id="loading">Loading...</div>
</body>
</html>
"@ | Out-File -Encoding utf8 build\index.html

Write-Host ""
Write-Host "Build complete! Files in ./build/" -ForegroundColor Green
Write-Host ""
Write-Host "To test locally:" -ForegroundColor Yellow
Write-Host "  cd build"
Write-Host "  python -m http.server 8080"
Write-Host "  Then open http://localhost:8080"
Write-Host ""
Write-Host "To deploy to Azure Static Web Apps:" -ForegroundColor Yellow
Write-Host "  az staticwebapp create -n ninesweeper -g <resource-group> --source ./build"