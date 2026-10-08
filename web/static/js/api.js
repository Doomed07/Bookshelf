// Обёртка над JSON API /api/v1.

export class ApiError extends Error {
  // detail — технический текст ошибки от сервера (поле "error"); по нему формы
  // отличают «имя занято» от «email занят». Пользователю его не показываем.
  constructor(status, message, detail = '') {
    super(message);
    this.status = status;
    this.detail = detail;
  }
}

// query собирает строку параметров, пропуская пустые значения: {title: null, limit: 37} → "?limit=37".
export function query(params) {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== null && value !== undefined && value !== '') search.set(key, value);
  }
  const s = search.toString();
  return s ? '?' + s : '';
}

export async function api(path, { method = 'GET', body } = {}) {
  let response;
  try {
    response = await fetch('/api/v1' + path, {
      method,
      headers: body === undefined ? {} : { 'Content-Type': 'application/json' },
      body: body === undefined ? undefined : JSON.stringify(body),
    });
  } catch {
    throw new ApiError(0, 'network error');
  }

  const text = await response.text();
  let data = null;
  try {
    data = text ? JSON.parse(text) : null;
  } catch {
    data = null;
  }

  if (!response.ok) {
    throw new ApiError(response.status, (data && data.message) || response.statusText, (data && data.error) || '');
  }
  return data;
}

// ---------- Аккаунт ----------
// Сессия живёт в HttpOnly-cookie: JS её не видит, браузер сам прикладывает её к запросам.

export const auth = {
  register: (username, email, password) =>
    api('/auth/register', { method: 'POST', body: { username, email, password } }),

  login: (login, password) => api('/auth/login', { method: 'POST', body: { login, password } }),

  // Тело {} нужно ради заголовка Content-Type: application/json: без него сервер
  // отклонит запрос (защита от подделки запросов с чужих сайтов).
  logout: () => api('/auth/logout', { method: 'POST', body: {} }),

  // me — вошедший пользователь или null, если не вошёл (401 здесь не ошибка).
  async me() {
    try {
      return await api('/auth/me');
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) return null;
      throw err;
    }
  },
};
