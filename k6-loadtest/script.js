import http from 'k6/http';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export default function () {
  let headers = {
    headers: { 'Content-Type': 'application/json' },
  };
  let order = {
    customerName: randomString(8)
  };

  let orderRes = http.post(`http://${__ENV.ORDER_HOST}:${__ENV.ORDER_PORT}/v1/order`, JSON.stringify(order), headers);

  for (var i = 0; i < 10; i++) {
    let item = {
      itemName: randomString(8),
      quantity: randomIntBetween(1, 50),
      price: randomIntBetween(1, 10000) / 10
    }
    http.put(`http://${__ENV.ORDER_HOST}:${__ENV.ORDER_PORT}/v1/order/${orderRes.json().orderID}/item`, JSON.stringify(item), headers);
  }
}
