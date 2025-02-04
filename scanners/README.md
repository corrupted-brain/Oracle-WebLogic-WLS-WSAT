# CVE-2017-10271 Vulnerability Scanner


Weblogic wls-wsat Component Deserialization Vulnerability (CVE-2017-10271) Detection Executable

### Usage

```
~/g/p/C/scanners ❯❯❯ ./bin/CVE-2017-10271.release.1.5.0.amd64.darwin -h
A purpose built scanner for detecting CVE-2017-10271. Starts a web
server on the LPORT and then logs any host which contacts it, as they are
vulnerable.

Example usage:
./CVE-2017-10271.release.1.5.0.amd64.linux -s "10.10.10.10" -t "$(pwd)/targets.txt -o output_file.txt -v --all-urls"

Example targets.txt:
http://pwned.com:7001/
https://pwnedalso.com:8002/

Usage:
  cve-2017-10271 [flags]
  cve-2017-10271 [command]

Available Commands:
  help        Help about any command
  version     The version of the binary

Flags:
  -u, --all-urls                Check for all possible vulnerable URL suffixes
  -h, --help                    help for cve-2017-10271
  -s, --listening-host string   The IP of this machine's public interface
  -l, --listening-port int      The port to listen for vulnerable responses (default 4444)
  -o, --output-file string      File to output results to
  -t, --target-file string      File with list of targets in http(s)://HOSTNAME:PORT format
  -a, --threads int             Number of threads to use while scanning (default 10)
  -v, --verbose                 Enable verbose mode (Print who is being scanned
  -w, --wait-time int           Seconds to wait after we complete sending payloads (default 20)

Use "cve-2017-10271 [command] --help" for more information about a command.
```
