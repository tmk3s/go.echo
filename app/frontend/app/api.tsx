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

  // セッション切れ（401）の場合はトップページへ戻す
  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (axios.isAxiosError(error) && error.response?.status === 401) {
        window.location.href = '/';
      }
      return Promise.reject(error);
    }
  )

  return instance;
}
