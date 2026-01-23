// load the api server with request

// https://k6.io/
import { post } from "k6/http";
import { check } from "k6";

export const options = {
  iterations: 1000,
};

export default function () {
  const payload = JSON.stringify({
    email: "test@example.com",
    balance: 1.1,
    age: 25,
    name: "John Doe",
  });
  const res = post("http://localhost:3000/users", payload);
  check(res, { "status was 200": (r) => r.status == 200 });
}

// load();
