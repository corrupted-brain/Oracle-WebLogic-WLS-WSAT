#!/bin/sh

echo "[+] Installing requests dependency"
pip3 install -U requests

echo "[+] Starting listener on port 4444"
python3 -m http.server 4444
