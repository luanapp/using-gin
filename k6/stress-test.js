import http from 'k6/http';
import { check } from 'k6';

export const options = {
    thresholds: {
        http_req_failed: [{threshold: 'rate<0.01', abortOnFail: true}],
        http_req_duration: ['p(99)<100'],
    },

    scenarios: {
        // arbitrary name of scenario
        average_load: {
            executor: 'ramping-vus',
            stages: [
                // ramp up to average load of 20 virtual users
                { duration: '10s', target: 20 },
                // maintain load
                { duration: '50s', target: 20 },
                // ramp down to zero
                { duration: '5s', target: 0 },
            ],
        },
    },
};

export default function () {
    const baseUrl = 'http://localhost:8080'
    const speciesUrl = `${baseUrl}/species`

    const payload = JSON.stringify({
        "scientific_name": "scientific_name",
        "genus": "genus",
        "family": "family",
        "order": "order",
        "class": "class",
        "phylum": "phylum",
        "kingdom": "kingdom"
    });
    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.post(speciesUrl, payload, params);

    check(res, {
        'response code was 201': (res) => res.status === 201,
    });

    const saved = res.json();
    res = http.get(speciesUrl+`/${saved.id}`);

    check(res, {
        'response code was 200': (res) => res.status === 200,
    });

    res = http.del(speciesUrl+`/${saved.id}`);
    check(res, {
        'response code was 202': (res) => res.status === 202,
    });
}