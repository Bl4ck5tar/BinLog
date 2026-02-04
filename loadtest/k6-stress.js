import http from 'k6/http'
import { check, sleep } from 'k6'

const API_BASE = __ENV.API_BASE || 'http://localhost:8080/api'

export const options = {
  stages: [
    { duration: '30s', target: 50 },
    { duration: '30s', target: 150 },
    { duration: '30s', target: 300 },
    { duration: '30s', target: 500 },
    { duration: '30s', target: 700 },
    { duration: '30s', target: 900 },
    { duration: '30s', target: 1100 },
    { duration: '30s', target: 900 },
    { duration: '30s', target: 500 },
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(50)<400', 'p(95)<1500', 'p(99)<3000'],
  },
}

function get(path) {
  return http.get(`${API_BASE}${path}`)
}

export default function () {
  const endpoints = [
    '/article/search?order=desc&sort=time&page=1&page_size=10',
    '/article/category',
    '/article/tags',
    '/website/info',
    '/website/news?source=zhihu',
    '/website/calendar',
    '/friendLink/info',
  ]

  for (const path of endpoints) {
    const res = get(path)
    check(res, {
      'status 200': (r) => r.status === 200,
      'code 0': (r) => {
        try {
          const data = r.json()
          return data && data.code === 0
        } catch (_) {
          return false
        }
      },
    })
  }

  sleep(1)
}
