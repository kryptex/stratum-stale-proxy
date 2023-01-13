### Usage

mine-with-stales.bat:
```
@cd /d "%~dp0"
start /b proxy.exe -local 127.0.0.1:8888 -upstream ethw.kryptex.network:7777 -delay 60000
t-rex.exe -a ethash -o stratum+tcp://127.0.0.1:8888 -u WALLET.delay60 -p x --send-stales
pause
```