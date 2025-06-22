# Ping Scanner Tool v0.1.0
Simple ICMP ping scanner written in Go. It reads ip addresses from a file, sends a ping to each address and writes the address in another file if it was successful.

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
Sample output:
```
2025/06/23 01:35:33 Read 3 IP addresses.      
2025/06/23 01:35:33 Starting ICMP ping scan...
2025/06/23 01:35:33 FAILURE: 10.6.0.1 did not respond
2025/06/23 01:35:33 FAILURE: 10.6.0.2 did not respond
2025/06/23 01:35:33 No response from any IP
```

### Multi-threaded
The multi-threading feature enables faster scan operations by equally dividing the ip addresses among the threads.
```
ping-scanner --threads <number of threads>
```
