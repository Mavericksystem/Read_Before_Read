# Benchmarks

## NIM admission and rate limiting

Configuration:

- `nimMaxInFlight`: 1
- Token bucket: 40 requests/minute, burst 1
- k6: 5 VUs for 90 seconds
- Endpoint: `POST /api/v1/analyze` with `https://example.com`

Run on 2026-09-10:

- Requests: 5
- Successful requests: 0
- Failed requests: 5 (100%)
- Average duration: 59.99s
- Median duration: 59.99s
- Throughput: 0.041666 requests/s
- Server behavior: one NIM request reached the external API; four requests waited for the throttle slot and were canceled when the 60s client timeout expired.

Finding: the Go limiter and single-flight admission are active, but this run cannot establish the intended 40/minute throughput because the NVIDIA NIM dependency did not complete within the k6 HTTP client's 60-second timeout. A successful capped-throughput benchmark requires a responsive NIM endpoint and should be rerun before claiming zero errors or approximately 60 successful requests in 90 seconds.
