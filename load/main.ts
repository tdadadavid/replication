// load the api server with request

// https://k6.io/
import { post } from "k6/http";
import { check, sleep } from "k6";

export const options = {
  iterations: 10000,
  // stages: [{ duration: "5s", target: 5 }],
};

export default function () {
  const id = __VU * 10000 + __ITER;

  const payload = JSON.stringify({
    email: "email@gmail.com" + id,
    balance: id,
    age: 25 + id,
    name: "David" + id,
  });

  const params = { headers: { "Content-Type": "application/json" } };
  const res = post("http://127.0.0.1:8000/async-users", payload, params);

  check(res, { "status was 200": (r) => r.status == 200 });

  sleep(0.01);
}
