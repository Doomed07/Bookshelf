// Общие куски интерфейса: шапка, подвал, обложки, оценки, пагинация, ошибки.
//
// Разметка собирается шаблоном html`…`: всё, что подставляется в ${…}, экранируется,
// кроме вложенных html`…`. Поэтому текст из базы (рецензии, имена) не может стать HTML-кодом.

import { ApiError } from './api.js';
import { avatarClass, coverClass, initial, scoreClass } from './format.js';

class SafeHTML {
  constructor(value) {
    this.value = value;
  }
}

export function escapeHTML(s) {
  return String(s)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

function renderValue(v) {
  if (v === null || v === undefined || v === false) return '';
  if (v instanceof SafeHTML) return v.value;
  if (Array.isArray(v)) return v.map(renderValue).join('');
  return escapeHTML(v);
}

export function html(strings, ...values) {
  let out = strings[0];
  values.forEach((v, i) => {
    out += renderValue(v) + strings[i + 1];
  });
  return new SafeHTML(out);
}

export function mount(element, content) {
  element.innerHTML = renderValue(content);
}

export const $ = (selector) => document.querySelector(selector);

// ---------- Шапка и подвал ----------

const NAV = [
  { key: 'signup', href: '/signup', label: 'Регистрация' },
  { key: 'books', href: '/books', label: 'Книги' },
  { key: 'users', href: '/users', label: 'Читатели' },
  { key: 'stats', href: '/stats', label: 'Статистика' },
];

function header({ nav = '', navDetail = false, search = '' }) {
  const links = NAV.map((item) => {
    const current = item.key === nav ? html` aria-current="${navDetail ? 'true' : 'page'}"` : '';
    return html`<a href="${item.href}"${current}>${item.label}</a>`;
  });

  return html`
    <header class="site-header">
      <div class="container site-header__inner">
        <a class="logo" href="/">
          <svg class="logo__mark" viewBox="0 0 22 20" aria-hidden="true" focusable="false">
            <rect x="0" y="5" width="6" height="15" rx="1" fill="#c69b5d"/>
            <rect x="8" y="0" width="6" height="20" rx="1" fill="#756ab4"/>
            <rect x="16" y="7" width="6" height="13" rx="1" fill="#b55280"/>
          </svg>
          <span class="logo__word">Shelfmate</span>
        </a>
        <nav class="main-nav" aria-label="Основное меню">${links}</nav>
        <form class="search" action="/books" method="get" role="search">
          <label class="visually-hidden" for="search-title">Поиск по названию книги</label>
          <input class="search__input" id="search-title" type="search" name="title" value="${search}" placeholder="Найти книгу…" maxlength="200" autocomplete="off">
          <button class="search__button" type="submit">
            <svg viewBox="0 0 24 24" aria-hidden="true" focusable="false"><circle cx="11" cy="11" r="7"/><path d="m20 20-4.2-4.2"/></svg>
            <span class="visually-hidden">Искать</span>
          </button>
        </form>
      </div>
    </header>`;
}

function footer() {
  return html`
    <footer class="site-footer">
      <div class="container site-footer__inner">
        <nav class="footer-nav" aria-label="Ссылки в подвале">
          <a href="/about">О проекте</a>
          <a href="/books">Книги</a>
          <a href="/users">Читатели</a>
          <a href="/stats">Статистика</a>
          <a href="/swagger/">API</a>
        </nav>
        <p>© 2026 Shelfmate. Сделано читателями для читателей.</p>
      </div>
    </footer>`;
}

// renderLayout вставляет шапку и подвал в #site-header и #site-footer страницы.
export function renderLayout(options = {}) {
  mount($('#site-header'), header(options));
  mount($('#site-footer'), footer());
}

export function setTitle(title) {
  document.title = title ? `${title} · Shelfmate` : 'Shelfmate — социальная сеть для тех, кто любит читать';
}

// ---------- Книги и оценки ----------

export function coverLabel(book) {
  return html`<span class="cover__title">${book.title}</span><span class="cover__author">${book.author}</span>`;
}

// cover — обложка-ссылка. extra — дополнительные классы, decorative — скрыть от скринридеров
// (рядом уже есть текстовая ссылка на ту же книгу).
export function cover(book, { small = false, extra = '', decorative = false, title = false } = {}) {
  const cls = ['cover', small ? 'cover--small' : '', coverClass(book.id), extra].filter(Boolean).join(' ');
  const hidden = decorative ? html` tabindex="-1" aria-hidden="true"` : '';
  const tip = title ? html` title="${book.title} — ${book.author}"` : '';
  return html`<a class="${cls}" href="/books/${book.id}"${hidden}${tip}>${coverLabel(book)}</a>`;
}

// score — значок оценки 1–100; для null ничего не выводит.
export function score(value) {
  if (value === null || value === undefined) return '';
  return html`<span class="score ${scoreClass(value)}"><span class="visually-hidden">Оценка </span>${value}</span>`;
}

export function avatar(userId, username, extra = '') {
  return html`<span class="avatar ${extra} ${avatarClass(userId)}" aria-hidden="true">${initial(username)}</span>`;
}

export function emptyState(content) {
  return html`<p class="empty-state">${content}</p>`;
}

export function loading() {
  return html`<p class="loading" role="status">Загрузка…</p>`;
}

// ---------- Параметры адреса и пагинация ----------

export const MAX_PAGE = 10000;

// pageParam — номер страницы из ?page=. Неверное значение — ошибка 400.
export function pageParam() {
  const raw = new URLSearchParams(location.search).get('page');
  if (raw === null || raw === '') return 1;
  const page = Number(raw);
  if (!Number.isInteger(page) || page < 1 || page > MAX_PAGE) throw new ApiError(400, 'invalid page');
  return page;
}

// pathID — id из адреса вида /books/3. Неверный id — «нет такой страницы» (404).
export function pathID() {
  const raw = location.pathname.split('/').filter(Boolean).pop();
  const id = Number(raw);
  if (!/^\d+$/.test(raw || '') || id < 1 || id > 2147483647) throw new ApiError(404, 'invalid id');
  return id;
}

// trimPage — приём «limit+1»: просим у API на один элемент больше размера страницы.
// Пришёл лишний — значит, есть следующая страница; сам он не показывается.
export function trimPage(items, size) {
  const list = items || [];
  return list.length > size ? [list.slice(0, size), true] : [list, false];
}

// readBooks оставляет книги, которые хоть раз прочитали.
export function readBooks(books, limit) {
  return (books || []).filter((b) => b.reads_count > 0).slice(0, limit);
}

// pageURL меняет только номер страницы, остальные параметры (title, read) сохраняет.
export function pageURL(page) {
  const params = new URLSearchParams(location.search);
  if (page <= 1) params.delete('page');
  else params.set('page', page);
  params.sort();
  const s = params.toString();
  return location.pathname + (s ? '?' + s : '');
}

export function pagination(page, hasNext) {
  if (page <= 1 && !hasNext) return '';
  return html`
    <nav class="pagination" aria-label="Страницы">
      ${page > 1 ? html`<a class="pagination__link" href="${pageURL(page - 1)}" rel="prev">← Назад</a>` : ''}
      <span class="pagination__current">Страница ${page}</span>
      ${hasNext ? html`<a class="pagination__link" href="${pageURL(page + 1)}" rel="next">Дальше →</a>` : ''}
    </nav>`;
}

// ---------- Ошибки ----------

const ERRORS = {
  400: ['Некорректный запрос', 'Проверьте адрес страницы или параметры фильтра.'],
  404: ['Страница не найдена', 'Возможно, её удалили или в адресе опечатка.'],
  409: ['Конфликт данных', 'Данные изменились, пока вы их смотрели. Обновите страницу.'],
};

export function errorView(status) {
  if (!document.querySelector('meta[name="robots"]')) {
    const meta = document.createElement('meta');
    meta.name = 'robots';
    meta.content = 'noindex';
    document.head.append(meta);
  }

  const code = ERRORS[status] ? status : 500;
  const [heading, message] = ERRORS[code] || ['Что-то пошло не так', 'Попробуйте обновить страницу чуть позже.'];
  setTitle(heading);
  return html`
    <div class="container page">
      <section class="error-page">
        <p class="error-page__code">${code}</p>
        <h1 class="error-page__title">${heading}</h1>
        <p class="error-page__text">${message}</p>
        <p class="error-page__links">
          <a class="button button--cta" href="/">На главную</a>
          <a href="/books">Каталог книг</a>
        </p>
      </section>
    </div>`;
}

// start рисует шапку и подвал и запускает загрузку страницы. Любая ошибка
// (404 из API, неверный ?page=, сбой сети) превращается в страницу ошибки в #app.
export async function start(layout, load) {
  renderLayout(layout);
  try {
    await load();
  } catch (err) {
    if (!(err instanceof ApiError)) console.error(err);
    mount($('#app'), errorView(err instanceof ApiError ? err.status : 500));
  }
}
