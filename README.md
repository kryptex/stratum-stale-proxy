### Build

1. Install Go: https://go.dev
2. `go build proxy.go`

### Usage
1. Download [T-Rex miner](https://github.com/trexminer/T-Rex/releases) (tolerates long delays in share submmitting)
2. Extract
3. Drop in proxy.exe
4. Use `mine-with-stales.bat` given below

**mine-with-stales.bat**:
```
@cd /d "%~dp0"
start /b proxy.exe -local 127.0.0.1:8888 -upstream ethw.kryptex.network:7777 -delay 60000
t-rex.exe -a ethash -o stratum+tcp://127.0.0.1:8888 -u WALLET.delay60 -p x --send-stales
pause
```
