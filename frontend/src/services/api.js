import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080/api", // Update with your backend URL
});

export default api;
