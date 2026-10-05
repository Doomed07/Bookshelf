import { api, ApiError, query } from '../api.js';
import { $, avatar, emptyState, html, mount, pageParam, pagination, setTitle, start, trimPage } from '../ui.js';

const PAGE_SIZE = 30;

start({ nav: 'users' }, async () => {
  const page = pageParam();
  const found = await api('/users' + query({ limit: PAGE_SIZE + 1, offset: (page - 1) * PAGE_SIZE }));

  const [users, hasNext] = trimPage(found, PAGE_SIZE);
  if (users.length === 0 && page > 1) throw new ApiError(404, 'page after the last');

  setTitle('Читатели');

  // API отдаёт и email, но это личные данные — на странице показываем только имя
  const list = users.length === 0
    ? emptyState(html`Читателей пока нет. <a href="/signup">Станьте первым</a>.`)
    : html`<ul class="user-grid">${users.map((u) => html`
        <li>
          <a class="user-card" href="/users/${u.id}">
            ${avatar(u.id, u.username, 'avatar--medium')}
            <span class="user-card__name">${u.username}</span>
          </a>
        </li>`)}</ul>`;

  mount($('#app'), html`
    <div class="container page">
      <header class="page-header">
        <h1 class="page-title">Читатели</h1>
        <p class="page-subtitle">В порядке регистрации.</p>
      </header>
      ${list}
      ${pagination(page, hasNext)}
    </div>`);
});
