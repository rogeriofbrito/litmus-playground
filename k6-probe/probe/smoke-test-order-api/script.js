import http from 'k6/http';
import { check, fail } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
    thresholds: {
        checks: [
            {
                threshold: 'rate>=1',
                abortOnFail: true,
            },
        ],
    },
};

export default function () {
    const headers = {
        headers: { 'Content-Type': 'application/json' },
    };
    const order = {
        customerName: randomString(8)
    };

    const orderRes = http.post(`http://${__ENV.ORDER_HOST}:${__ENV.ORDER_PORT}/v1/order`, JSON.stringify(order), headers);
    if (!check(orderRes, { 'order create status code MUST be 200': (res) => res.status == 200, })) {
        fail('order create status code was *not* 200');
    }

    const item = {
        itemName: randomString(8),
        quantity: randomIntBetween(1, 50),
        price: randomIntBetween(1, 10000) / 10
    }

    const itemRes = http.put(`http://${__ENV.ORDER_HOST}:${__ENV.ORDER_PORT}/v1/order/${orderRes.json().orderID}/item`, JSON.stringify(item), headers);
    if (!check(itemRes, { 'item create status code MUST be 200': (res) => res.status == 200, })) {
        fail('item create status code was *not* 200');
    }
}
