import http from 'k6/http';
export const options = { vus: 5, duration: '90s' };
export default function () {
    http.post('http://localhost:8080/api/v1/analyze',
        JSON.stringify({ url: "https://example.com" }),
        { headers: { 'Content-Type': 'application/json' } });
}