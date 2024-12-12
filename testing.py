# import requests
# import json
# from concurrent.futures import ThreadPoolExecutor

# # Define the target URL and payload
# url = "http://localhost:3000/api/v1/compile"
# payload = {
#     "code": "for i in range(1, 100):\n    print(i)",
#     "lang": "py",
#     "input": ""
# }
# headers = {
#     "Content-Type": "application/json"
# }

# # Function to send a single request
# def send_request():
#     try:
#         response = requests.post(url, headers=headers, data=json.dumps(payload))
#         print(f"Status Code: {response.status_code}, Response: {response.text}")
#     except requests.exceptions.RequestException as e:
#         print(f"Request failed: {e}")

# # Number of concurrent requests
# num_requests = 30

# # Use ThreadPoolExecutor for concurrent requests
# with ThreadPoolExecutor(max_workers=num_requests) as executor:
#     futures = [executor.submit(send_request) for _ in range(num_requests)]

# # Ensure all requests are completed
# for future in futures:
#     try:
#         future.result()
#     except Exception as e:
#         print(f"Error during request execution: {e}")

# import requests
# import concurrent.futures
# import time

# # Define the URL and headers
# url = "https://codeploy-f46rfl3ysa-el.a.run.app/api/compiler"

# headers = {
#     "accept": "*/*",
#     "accept-language": "en-US,en;q=0.9,hi;q=0.8",
#     "content-type": "application/json",
#     "origin": "https://playground.compile-me.com",
#     "priority": "u=1, i",
#     "referer": "https://playground.compile-me.com/",
#     "sec-ch-ua": "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
#     "sec-ch-ua-mobile": "?0",
#     "sec-ch-ua-platform": "\"Windows\"",
#     "sec-fetch-dest": "empty",
#     "sec-fetch-mode": "cors",
#     "sec-fetch-site": "cross-site",
#     "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
# }

# data = {
#     "code": "// Write your code here\nconsole.log(100)",
#     "language": "js",
#     "stdInput": ""
# }

# # Function to send a request and print the response
# def send_request():
#     try:
#         response = requests.post(url, headers=headers, json=data)
#         # Print the status code and response body
#         # print(f"Response Status: {response.status_code}")
#         print(f"Response Body: {response.text}")  # Print raw response body
#     except requests.exceptions.RequestException as e:
#         print(f"Request failed: {e}")

# # Main function to run the load test
# def load_test(num_requests_per_second, duration_seconds):
#     start_time = time.time()

#     # Create a thread pool to send requests concurrently
#     with concurrent.futures.ThreadPoolExecutor(max_workers=num_requests_per_second) as executor:
#         while time.time() - start_time < duration_seconds:
#             # Schedule the requests to run concurrently
#             futures = [executor.submit(send_request) for _ in range(num_requests_per_second)]
#             concurrent.futures.wait(futures)

#             # Optional: Delay to control the rate of requests
#             time.sleep(1)  # This ensures that we send 50 requests per second

# # Run the load test for 50 requests per second, for 10 seconds
# load_test(5000, 10)

import requests
import concurrent.futures
import time

# Define the URL and headers
url = "https://codeploy-f46rfl3ysa-el.a.run.app/api/compiler"

headers = {
    "accept": "*/*",
    "accept-language": "en-US,en;q=0.9,hi;q=0.8",
    "content-type": "application/json",
    "origin": "https://playground.compile-me.com",
}

data = {
    "code": "// Function to check if a number is prime\nfunction isPrime(num) {\n    if (num <= 1) return false; // 0 and 1 are not prime\n    for (let i = 2; i <= Math.sqrt(num); i++) {\n        if (num % i === 0) return false; // Divisible by any number other than 1 and itself\n    }\n    return true; // It's a prime number\n}\n\n// Function to count prime numbers up to a limit\nfunction countPrimes(limit) {\n    let primeCount = 0;\n    for (let i = 2; i <= limit; i++) {\n        if (isPrime(i)) {\n            primeCount++;\n        }\n    }\n    return primeCount;\n}\n\n// Counting primes from 1 to 1000000\nconst primeCount = countPrimes(1000000);\nconsole.log(primeCount);",
    "language": "js",
    "stdInput": ""
}

# Function to send a single request and print the response
def send_request():
    try:
        response = requests.post(url, headers=headers, json=data)
        # Print the status code and response body
        # print(f"Response Status: {response.status_code}")
        print(f"Response Body: {response.text}")  # Print raw response body
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")

# Main function to run the load test with a given number of concurrent requests
def load_test(num_requests):
    # Start time for the load test
    start_time = time.time()

    # Use ThreadPoolExecutor to send requests concurrently
    with concurrent.futures.ThreadPoolExecutor(max_workers=num_requests) as executor:
        futures = [executor.submit(send_request) for _ in range(num_requests)]
        # Wait for all the futures to complete
        concurrent.futures.wait(futures)

    # Calculate the time taken for the load test
    end_time = time.time()
    print(f"Load test completed in {end_time - start_time:.2f} seconds.")

# Run the load test with 5000 concurrent requests
load_test(5000)
