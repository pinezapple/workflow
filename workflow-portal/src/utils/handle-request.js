import axios from "axios";

const service = axios.create({
  baseURL: "https://workflow.com",
  timeout: 5000,
  // withCredentials: true,
  headers: {
    Authorization:
      "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6ImZlbmNlX2tleV9rZXlzIn0.eyJwdXIiOiJhY2Nlc3MiLCJhdWQiOlsib3BlbmlkIiwidXNlciIsImNyZWRlbnRpYWxzIiwiZGF0YSIsImFkbWluIiwiZ29vZ2xlX2NyZWRlbnRpYWxzIiwiZ29vZ2xlX3NlcnZpY2VfYWNjb3VudCIsImdvb2dsZV9saW5rIiwiZ2E0Z2hfcGFzc3BvcnRfdjEiXSwic3ViIjoiMzMxIiwiaXNzIjoiaHR0cHM6Ly9nZW5vbWUudmluYmlnZGF0YS5vcmcvdXNlciIsImlhdCI6MTYyMjkwNDE3NCwiZXhwIjoxNjIyOTQwMTc0LCJqdGkiOiIxMTliMWE1MC1hODJkLTRmODgtOTQwMi1mNTQwODZjMmNlYzAiLCJzY29wZSI6WyJvcGVuaWQiLCJ1c2VyIiwiY3JlZGVudGlhbHMiLCJkYXRhIiwiYWRtaW4iLCJnb29nbGVfY3JlZGVudGlhbHMiLCJnb29nbGVfc2VydmljZV9hY2NvdW50IiwiZ29vZ2xlX2xpbmsiLCJnYTRnaF9wYXNzcG9ydF92MSJdLCJjb250ZXh0Ijp7InVzZXIiOnsibmFtZSI6Im5ndXllbnF1YW4udHJhZGVAZ21haWwuY29tIiwiaXNfYWRtaW4iOmZhbHNlLCJnb29nbGUiOnsicHJveHlfZ3JvdXAiOm51bGx9LCJwcm9qZWN0cyI6e319fSwiYXpwIjoiIn0.HyeNjSEp_Hf9ahMihSuKPMW7FLcAJsuU1fDgvvoW2Ub1Kpk5D3EV_eFz58-jLU1Z5xLkKAAQBy75ow8JNsmAw3QidkuoVwdeMVltpTFsroMDM2nCIhkKFT1ZZqsa1t0DImn1cVyn9B2wOCW7xlH6xJV4mYSdBtTzyuL2XDnFW8ya1G29Jr3ekgph5Tz3SRGdjckYzco6mZWAAY_c7E5fjuAMRnhQhaP2LO2l4pWmWHVojdCWmQVzQGZQTWLQ-QRbWAS-6Yl-7CCfMCM5PvVQMtMo0MGk7PLGvxSnleSW_H6MKDHDrHwBPiCy76zgteZ4MK7WhYzCd7l7Ah_-u1sN9A",
  },
});

service.interceptors.request.use(
  (config) => {
    return config;
  },
  (error) => {
    console.log("Error when request: " + error);
    return Promise.reject(error);
  }
);

service.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.log("Error in response: " + error.message);
    return Promise.reject(error);
  }
);

export default service;
