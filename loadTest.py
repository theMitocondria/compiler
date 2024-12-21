import requests
import concurrent.futures
import time

# Define the server endpoint URL
URL = "http://localhost:3000/api/v1/compile"

# Define the payload to send with each request
PAYLOAD ={
  "lang": "cpp",
  "code": "#include <bits/stdc++.h> \n using namespace std;\n#include <vector>\n\nstd::vector<int> sieve(int n) {\n    std::vector<bool> isPrime(n + 1, true);\n    isPrime[0] = isPrime[1] = false;\n\n    for (int i = 2; i * i <= n; i++) {\n        if (isPrime[i]) {\n            for (int j = i * i; j <= n; j += i) {\n                isPrime[j] = false;\n            }\n        }\n    }\n\n    std::vector<int> primes;\n    for (int i = 2; i <= n; i++) {\n        if (isPrime[i]) {\n            primes.push_back(i);\n        }\n    }\n    return primes;\n}\n\nint main() {\n    int n;cin>>n;\n    std::vector<int> primes = sieve(n);\n    for (size_t i = 0; i < primes.size(); ++i) {\n        if (i > 0) std::cout << \", \";\n        std::cout << primes[i];\n    }\n    std::cout << std::endl;\n    return 0;\n}",
  "input": "20"
}

# Function to send a request to the server
def send_request():
    try:
        response = requests.post(URL, json=PAYLOAD, timeout=100)
        return response.status_code, response.text
    except requests.exceptions.RequestException as e:
        return "Error", str(e)


for i in range(1000):
    print(send_request())
