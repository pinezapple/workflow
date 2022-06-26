import axios from "axios";

const service = axios.create({
  baseURL: "http://localhost:8084",
  timeout: 5000,
  // withCredentials: true,
  headers: {
    Authorization:
      "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImZlbmNlX2tleV9rZXlzIn0.eyJwdXIiOiJhY2Nlc3MiLCJhdWQiOlsib3BlbmlkIiwidXNlciIsImNyZWRlbnRpYWxzIiwiZGF0YSIsImFkbWluIiwiZ29vZ2xlX2NyZWRlbnRpYWxzIiwiZ29vZ2xlX3NlcnZpY2VfYWNjb3VudCIsImdvb2dsZV9saW5rIiwiZ2E0Z2hfcGFzc3BvcnRfdjEiXSwic3ViIjoiNDEiLCJpc3MiOiJodHRwczovL2dlbm9tZS52aW5iaWdkYXRhLm9yZy91c2VyIiwiaWF0IjoxNjI4NjQ4NjUwLCJleHAiOjE2Mjg2ODQ2NTAsImp0aSI6IjY5NmU2MjM1LThlMjAtNDJhOC05NTQwLWFjNGU5OTUxMDRlNyIsInNjb3BlIjpbIm9wZW5pZCIsInVzZXIiLCJjcmVkZW50aWFscyIsImRhdGEiLCJhZG1pbiIsImdvb2dsZV9jcmVkZW50aWFscyIsImdvb2dsZV9zZXJ2aWNlX2FjY291bnQiLCJnb29nbGVfbGluayIsImdhNGdoX3Bhc3Nwb3J0X3YxIl0sImNvbnRleHQiOnsidXNlciI6eyJuYW1lIjoiZG9uZ29jdHVhbi4wMTAxQGdtYWlsLmNvbSIsImlzX2FkbWluIjp0cnVlLCJnb29nbGUiOnsicHJveHlfZ3JvdXAiOm51bGx9LCJwcm9qZWN0cyI6e319fSwiYXpwIjoiIn0.Occ6mavtYtiLM5AxCBK8rzAyvot7MojeGlptkJjYy5S_F8vyijxh6U6WWILbK4rUArbPvuDV0sopVJfThMUthtVAKiVUpUIMmxwGO60faZw3S7pKyvv95aP5i-HZO9mYbiEQXNEF4Q_sw0Q4CIUN9rM0wG9m8S7_XMNGUnJ_BqtuZBANY7WS7GvxTS_ECn-VBeFOX4lj1GTOCbsikmDsjDFLo7kpnnFXZvK5hDs8OCisP_Pgi8ajie7kpIxSFFihWzHypjNVCbRY6PGCwQjjzRYU9AO2lEc6Jb58Enz3VfFNDoY90C2x7FzdS3S9R6pMozl7Egi99b7reRK7AwCmxA",
  },
});

service.interceptors.request.use(
  (config) => {
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

service.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export default service;
