"use client";

import axios, {AxiosInstance} from 'axios';

export default (): AxiosInstance => {
  const instance = axios.create({
    baseURL: 'http://localhost:1323',
    headers: {
      'Content-Type': 'application/json',
      'X-Requested-With': 'XMLHttpRequest',
    },
    withCredentials: true,
    responseType: 'json'
  })
  return instance;
}
