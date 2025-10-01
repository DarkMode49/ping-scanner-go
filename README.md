# Ping Scanner Tool
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)  

ICMP ping scanner written in Go. It efficiently reads IP addresses from a file, sends a ping to each address and writes the address in an output file if it was responsive.

<details>
<summary>✅ IPv6 support</summary>
Write version 6 or 4 in input file<br/>
No problem should arise
</details>

<details>
<summary>✅ Multi-threading capability</summary>
A large input file with hundreds of thousands of lines?<br/>
Pass a thread number argument and the file will be equally divided across threads working simultaneously throught their part<br/>
</details>

<details>
<summary>✅ Memory efficient</summary>
Each thread will read one IP or line at a time. So NO loading heavy input file<br/>
Upon discovering a responsive IP, it will be written to the output file without interrupting anything else. So NO accumulating all the results in RAM<br/>
</details>

## Usage
### Basic
1. Create a text file with name `ips.txt` or any other name:
```
140.82.121.4
2a00:1450:4014:80e::200e
151.101.130.219
151.101.66.219
151.101.194.219
151.101.2.219
...
```
2. Run the executable binary:
```shell
ping-scanner.exe
```
```shell
# may require super user
./ping-scanner
```
Output
```
[2025-10-02 01:44:35][INFO] Ping Scanner v1.0.0       
[2025-10-02 01:44:35][WARN] Starting ICMP ping scan...
[2025-10-02 01:44:35][INFO] Pinging 140.82.121.4      
[2025-10-02 01:44:36][INFO] [Thread 0] 32 bytes from 140.82.121.4 icmp_seq=0 time=116.7838ms
[2025-10-02 01:44:36][INFO] Pinging 2a00:1450:4014:80e::200e
[2025-10-02 01:44:36][INFO] [Thread 0] 32 bytes from 2a00:1450:4014:80e::200e icmp_seq=0 time=172.0663ms
[2025-10-02 01:44:36][INFO] Pinging 151.101.130.219
[2025-10-02 01:44:36][INFO] [Thread 0] 32 bytes from 151.101.130.219 icmp_seq=0 time=132.342ms
[2025-10-02 01:44:36][INFO] Pinging 151.101.66.219
[2025-10-02 01:44:36][INFO] [Thread 0] 32 bytes from 151.101.66.219 icmp_seq=0 time=94.4346ms
[2025-10-02 01:44:36][INFO] Pinging 151.101.194.219
[2025-10-02 01:44:36][INFO] [Thread 0] 32 bytes from 151.101.194.219 icmp_seq=0 time=138.1868ms
[2025-10-02 01:44:36][INFO] Scan finished
[2025-10-02 01:44:36][INFO] Successful: 5  Failed: 0
```
Default `good.txt`:
```
140.82.121.4
2a00:1450:4014:80e::200e
151.101.130.219
151.101.66.219
151.101.194.219
```

### Multi-threaded
The multi-threading feature enables faster scan operations by equally dividing the ip addresses among the threads.  
Default: 1  
Recommended value: Number of CPU cores
```
ping-scanner --input large-file.txt --threads <number of threads>
```

## Download
Direct download link of latest release should be available at [release page](https://github.com/DarkMode49/ping-scanner-go/releases).

## Build
If for any reason you need to build it
```shell
go get .
go build .
```
The build command shall give you the compiled binary based on the operating system it ran upon.  
Set shell variables `GOOS` as for the operating system and `GOARCH` as for the CPU architecture you want to produce a build for.

## Change Log

### v1.0.0
Added IPv6 support  
As soon as a responsive IP was found, it will be written to the output file  
Improved log output  
SOLID compliance increased  
Silence argument now works properly  
Added max-errors parameter to feature error tolerance control

### v0.2.0
Improved console log output  
Replaced unreliable ping method  
Added silent argument  
Now input file is read directly instead of loading it into RAM  
