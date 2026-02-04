import http from 'k6/http'
import { check, sleep } from 'k6'

const API_BASE = __ENV.API_BASE || 'http://localhost:8080/api'

export const options = {
  stages: [
    { duration: '20s', target: 20 },
    { duration: '40s', target: 60 },
    { duration: '40s', target: 120 },
    { duration: '30s', target: 60 },
    { duration: '20s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(50)<200', 'p(95)<800', 'p(99)<1500'],
  },
}

function get(path) {
  return http.get(`${API_BASE}${path}`)
}

export default function () {
  const endpoints = [
    { name: 'article_search', path: '/article/search?order=desc&sort=time&page=1&page_size=10' },
    { name: 'article_category', path: '/article/category' },
    { name: 'article_tags', path: '/article/tags' },
    { name: 'website_info', path: '/website/info' },
    { name: 'website_news', path: '/website/news?source=zhihu' },
    { name: 'website_calendar', path: '/website/calendar' },
    { name: 'friend_link', path: '/friendLink/info' },
  ]

  for (const ep of endpoints) {
    const res = get(ep.path)
    check(res, {
      [`${ep.name} status 200`]: (r) => r.status === 200,
      [`${ep.name} success code`]: (r) => {
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
