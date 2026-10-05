import { api, ApiError, query } from '../api.js';
import { $, cover, emptyState, html, mount, pageParam, pagination, score, setTitle, start, trimPage } from '../ui.js';

const PAGE_SIZE = 36;

const search = (new URLSearchParams(location.search).get('title') || '').trim();

start({ nav: 'books', search }, async () => {
  const page = pageParam();
  const found = await api('/books' + query({
    title: search,
    limit: PAGE_SIZE + 1,
    offset: (page - 1) * PAGE_SIZE,
  }));

  const [books, hasNext] = trimPage(found, PAGE_SIZE);
  if (books.length === 0 && page > 1) throw new ApiError(404, 'page after the last');

  setTitle(search ? `Поиск: «${search}»` : 'Книги');

  const heading = search
    ? html`<h1 class="page-title">Поиск: «${search}»</h1>
           <p class="page-subtitle"><a href="/books">Сбросить поиск</a></p>`
    : html`<h1 class="page-title">Книги</h1>
           <p class="page-subtitle">Каталог по алфавиту. Чтобы найти книгу по названию, воспользуйтесь поиском в шапке.</p>`;

  const list = books.length === 0
    ? emptyState(search ? `По запросу «${search}» ничего не найдено.` : 'Каталог пуст.')
    : html`<ul class="poster-row">${books.map((b) => html`
        <li>
          ${cover(b)}
          <p class="poster-row__meta">${b.year}${b.score != null ? html` · ${score(b.score)}` : ''}</p>
        </li>`)}</ul>`;

  mount($('#app'), html`
    <div class="container page">
      <header class="page-header">${heading}</header>
      ${list}
      ${pagination(page, hasNext)}
    </div>`);
});
