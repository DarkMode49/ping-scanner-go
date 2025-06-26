# Ping Scanner Tool v0.2.0
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)  

Simple ICMP ping scanner written in Go. It reads ip addresses from a file, sends a ping to each address and writes the address in output file if it was successful.

✅ Multi-threading capability

## Usage
### Basic
1. Create a text file with name 'ips.txt':
```
10.6.0.1
10.6.0.2
10.6.0.3
...
```
2. Run the executable binary:
```shell
ping-scanner
```
Output
```
[2025-06-27 00:32:49] Ping Scanner v0.2.0
[2025-06-27 00:32:49] Starting ICMP ping scan...
[2025-06-27 00:32:49] Pinging 192.168.1.1
[2025-06-27 00:32:49] 32 bytes from 192.168.1.1: icmp_seq=0 time=7.6056ms
[2025-06-27 00:32:49] SUCCESS: 192.168.1.1
[2025-06-27 00:32:49] Pinging 192.168.1.2
[2025-06-27 00:32:51] FAILURE: 192.168.1.2
[2025-06-27 00:32:51] Pinging 192.168.1.3
[2025-06-27 00:32:53] FAILURE: 192.168.1.3
[2025-06-27 00:32:53] Pinging 192.168.1.4
[2025-06-27 00:32:55] FAILURE: 192.168.1.4
[2025-06-27 00:32:55] Pinging 192.168.1.5
[2025-06-27 00:32:57] FAILURE: 192.168.1.5
[2025-06-27 00:32:57] Wrote 1 responsive IP addresses to good.txt
```

### Multi-threaded
The multi-threading feature enables faster scan operations by equally dividing the ip addresses among the threads.  
Default: 1  
Recommended value: Number of CPU cores
```
ping-scanner --threads <number of threads>
```

## Download
Currently a pre-release version is available at the [release page](https://github.com/DarkMode49/ping-scanner-go/releases).

## Change Log

### v0.2.0
Improved console log output  
Replaced unreliable ping method  
Added silent argument  
Now input file is read directly instead of loading it into RAM  
