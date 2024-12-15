import requests
import json
import concurrent.futures

# Define the endpoint and headers
url = 'http://localhost:3000/api/v1/compile'
headers = {
    'Content-Type': 'application/json',
}

# Define the payload
payload = {
    "lang": "js",
    "input": "",
    "code": "for(var i = 0;i < 10000; i++){} console.log('done')"
}

# Function to send a POST request
def send_request():
    response = requests.post(url, headers=headers, data=json.dumps(payload))
    return response.status_code, response.text

# Send 100 requests concurrently
def main():
    with concurrent.futures.ThreadPoolExecutor(max_workers=100) as executor:
        futures = [executor.submit(send_request) for _ in range(10)]
        for future in concurrent.futures.as_completed(futures):
            status_code, response_text = future.result()
            print(f"Status Code: {status_code}, Response: {response_text}")
            print()

if __name__ == '__main__':
    main()
