import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 10 },
        { duration: '20s', target: 30 },
    ],
};

export default function () {
    // Step 1: Login
    const loginPayload = JSON.stringify({
        email: 'admin@vivasoftltd.com',
        password: '12345',
    });

    const loginHeaders = { 'Content-Type': 'application/json' };

    const loginRes = http.post('http://13.229.123.94:8080/v1/auth/login', loginPayload, {
        headers: loginHeaders,
    });

    check(loginRes, {
        'login succeeded': (r) => r.status === 200,
        'has access_token': (r) => r.json('access_token') !== undefined,
    });
    sleep(1);
}
