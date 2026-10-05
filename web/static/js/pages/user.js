import { api, ApiError, query } from '../api.js';
import { date, duration, isoDate, number, percent, plural } from '../format.js';
import {
  $, avatar, cover, emptyState, html, mount, pageParam, pagination, pathID, score, setTitle, start, trimPage,
} from '../ui.js';

const SHELF_PAGE_SIZE = 24;
const ACTIVITY_LIMIT = 10;

// tabFromQuery: ?read=true — прочитанные, ?read=false — «хочу прочитать», без параметра — все.
function tabFromQuery() {
  const read = new URLSearchParams(location.search).get('read');
  if (read === null) return { tab: 'all', read: null };
  if (read === 'true') return { tab: 'read', read: true };
  if (read === 'false') return { tab: 'unread', read: false };
  throw new ApiError(400, 'invalid read filter');
}

function tabs(current, stats) {
  const path = location.pathname;
  const items = [
    { key: 'all', label: 'Все', url: path, count: stats.books_on_shelf },
    { key: 'read', label: 'Прочитанные', url: path + '?read=true', count: stats.books_read },
    { key: 'unread', label: 'Хочу прочитать', url: path + '?read=false', count: stats.books_on_shelf - stats.books_read },
  ];
  return html`<nav class="tabs" aria-label="Полка">${items.map((t) => html`
    <a class="tabs__link" href="${t.url}"${t.key === current ? html` aria-current="page"` : ''}>${t.label} <span class="tabs__count">${number(t.count)}</span></a>`)}</nav>`;
}

function shelf(items, tab) {
  if (items.length === 0) {
    if (tab === 'read') return emptyState('Прочитанных книг пока нет.');
    if (tab === 'unread') return emptyState('Список «хочу прочитать» пуст.');
    return emptyState(html`Полка пока пуста. Книги на полку пока добавляются через <a href="/swagger/">API</a>.`);
  }
  return html`<ul class="poster-row poster-row--compact">${items.map(({ shelf: s, book }) => html`
    <li>
      ${cover(book)}
      <p class="poster-row__meta">${s.read ? (s.rating != null ? score(s.rating) : 'прочитано') : 'в планах'}</p>
    </li>`)}</ul>`;
}

function stat(value, label, extra = '') {
  return html`<div class="stat"><dt class="stat__label">${label}</dt><dd class="stat__value ${extra}">${value}</dd></div>`;
}

function activity(events) {
  if (events.length === 0) return emptyState('Активности пока нет.');
  return html`<ol class="activity">${events.map((e) => html`
    <li class="activity__item">
      <p class="activity__label">${e.name === 'finished' ? 'Прочитано' : 'Добавлено на полку'} · <time datetime="${isoDate(e.at)}">${date(e.at)}</time></p>
      <p class="activity__book"><a href="/books/${e.book_id}">${e.title}</a> ${score(e.rating)}</p>
      <p class="activity__author">${e.author}</p>
      ${e.review ? html`<p class="activity__review">${e.review}</p>` : ''}
    </li>`)}</ol>`;
}

start({ nav: 'users', navDetail: true }, async () => {
  const id = pathID();
  const { tab, read } = tabFromQuery();
  const page = pageParam();

  const [user, stats, found, events] = await Promise.all([
    api(`/users/${id}`),
    api('/stats' + query({ user_id: id })),
    api(`/users/${id}/bookshelf` + query({ read, limit: SHELF_PAGE_SIZE + 1, offset: (page - 1) * SHELF_PAGE_SIZE })),
    api(`/users/${id}/activity` + query({ limit: ACTIVITY_LIMIT })),
  ]);

  const [items, hasNext] = trimPage(found, SHELF_PAGE_SIZE);
  if (items.length === 0 && page > 1) throw new ApiError(404, 'page after the last');

  setTitle(user.username);

  mount($('#app'), html`
    <div class="container page">
      <header class="profile">
        ${avatar(user.id, user.username, 'avatar--large')}
        <h1 class="page-title">${user.username}</h1>
      </header>

      <div class="profile-layout">
        <section aria-labelledby="shelf-title">
          <h2 class="visually-hidden" id="shelf-title">Полка</h2>
          ${tabs(tab, stats)}
          ${shelf(items, tab)}
          ${pagination(page, hasNext)}
        </section>

        <aside class="profile-aside">
          <section aria-labelledby="user-stats-title">
            <div class="section-heading">
              <h2 class="section-heading__title" id="user-stats-title">Статистика</h2>
            </div>
            <dl class="stat-grid stat-grid--two">
              ${stat(number(stats.books_on_shelf), plural(stats.books_on_shelf, 'книга на полке', 'книги на полке', 'книг на полке'))}
              ${stat(number(stats.books_read), plural(stats.books_read, 'книга прочитана', 'книги прочитаны', 'книг прочитано'))}
              ${stat(percent(stats.read_percent), 'прочитано из добавленного')}
              ${stat(duration(stats.avg_read_time_hours), 'в среднем на книгу')}
              ${stat(stats.favorite_genre || '—', 'любимый жанр', 'stat__value--text')}
              ${stat(stats.favorite_author || '—', 'любимый автор', 'stat__value--text')}
              ${stat(number(stats.reviews_count), plural(stats.reviews_count, 'рецензия', 'рецензии', 'рецензий'))}
            </dl>
          </section>

          <section class="section" aria-labelledby="activity-title">
            <div class="section-heading">
              <h2 class="section-heading__title" id="activity-title">Активность</h2>
            </div>
            ${activity(events || [])}
          </section>
        </aside>
      </div>
    </div>`);
});
