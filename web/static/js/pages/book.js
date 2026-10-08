import { api, ApiError, query } from '../api.js';
import { coverClass, date, decimal, isoDate, number, plural, pluralCount } from '../format.js';
import {
  $, avatar, coverLabel, currentUser, emptyState, html, mount, pageParam, pagination, pathID, score, setTitle, start, trimPage,
} from '../ui.js';
import { bindShelfPanel, getShelfBook, shelfPanel } from '../shelf.js';

const REVIEWS_PAGE_SIZE = 10;

// ratingHistogram раскладывает оценки 1–100 по десяти диапазонам: 1–10, 11–20, …
function ratingHistogram(distribution) {
  const buckets = Array.from({ length: 10 }, (_, i) => ({ from: i * 10 + 1, to: (i + 1) * 10, count: 0, height: 0 }));
  let total = 0;
  for (const { rating, count } of distribution || []) {
    if (rating < 1 || rating > 100) continue;
    buckets[Math.floor((rating - 1) / 10)].count += count;
    total += count;
  }
  const maxCount = Math.max(...buckets.map((b) => b.count));
  for (const b of buckets) {
    // у маленького столбика минимальная высота, чтобы он не пропал
    if (b.count > 0) b.height = Math.max(Math.floor((b.count * 100) / maxCount), 4);
  }
  return [buckets, total];
}

function stat(value, label) {
  return html`<div class="stat"><dt class="stat__label">${label}</dt><dd class="stat__value">${value}</dd></div>`;
}

function reviewItem(r) {
  return html`
    <li class="review">
      <div class="review__body">
        <p class="review__meta">
          ${avatar(r.user_id, r.username)}
          <a class="review__author" href="/users/${r.user_id}">${r.username}</a>
          ${score(r.rating)}
        </p>
        ${r.review
          ? html`<p class="review__text review__text--full">${r.review}</p>`
          : html`<p class="review__text review__text--muted">Оценка без рецензии</p>`}
        ${r.reviewed_at
          ? html`<p class="review__date">Опубликовано <time datetime="${isoDate(r.reviewed_at)}">${date(r.reviewed_at)}</time></p>`
          : html`<p class="review__date">Прочитано <time datetime="${isoDate(r.read_at)}">${date(r.read_at)}</time></p>`}
      </div>
    </li>`;
}

// myShelfEntry — запись этой книги на полке вошедшего: объект, null (нет на полке)
// или undefined (не удалось загрузить — страница всё равно показывается).
async function myShelfEntry(me, bookId) {
  if (!me) return null;
  try {
    return await getShelfBook(me.id, bookId);
  } catch {
    return undefined;
  }
}

// load загружает и рисует страницу целиком. После любого действия с полкой страница
// перерисовывается заново: счётчики, гистограмма и список рецензий не устаревают.
async function load({ notice = '', focusPanel = false } = {}) {
  const id = pathID();
  const page = pageParam();
  const me = await currentUser();

  const [book, stats, found, entry] = await Promise.all([
    api(`/books/${id}`),
    api('/stats' + query({ book_id: id })),
    api(`/books/${id}/reviews` + query({ limit: REVIEWS_PAGE_SIZE + 1, offset: (page - 1) * REVIEWS_PAGE_SIZE })),
    myShelfEntry(me, id),
  ]);

  const [reviews, hasNext] = trimPage(found, REVIEWS_PAGE_SIZE);
  if (reviews.length === 0 && page > 1) throw new ApiError(404, 'page after the last');

  const [histogram, ratingsTotal] = ratingHistogram(stats.rating_distribution);
  setTitle(`${book.title} — ${book.author}`);

  mount($('#app'), html`
    <div class="container page">
      <article class="book">
        <div class="cover ${coverClass(book.id)} book__cover" aria-hidden="true">${coverLabel(book)}</div>
        <div class="book__info">
          <h1 class="book__title">${book.title}</h1>
          <p class="book__author">${book.author}</p>
          <p class="book__meta">${book.year} · ${pluralCount(book.pages, 'страница', 'страницы', 'страниц')}${book.score != null ? html` · ${score(book.score)}` : ''}</p>
          ${book.genres && book.genres.length
            ? html`<ul class="chips" aria-label="Жанры">${book.genres.map((g) => html`<li class="chip">${g}</li>`)}</ul>`
            : ''}
          <p class="book__description">${book.description}</p>
          <div id="shelf-panel" tabindex="-1">${shelfPanel({ me, entry, notice })}</div>
        </div>
      </article>

      <section class="section" aria-labelledby="book-stats-title">
        <div class="section-heading">
          <h2 class="section-heading__title" id="book-stats-title">Книга в цифрах</h2>
        </div>
        <dl class="stat-grid">
          ${stat(number(stats.users_on_shelf), 'на полках')}
          ${stat(number(stats.users_read), plural(stats.users_read, 'прочтение', 'прочтения', 'прочтений'))}
          ${stat(decimal(stats.average_rating), 'средняя оценка')}
          ${stat(number(stats.reviews_count), plural(stats.reviews_count, 'рецензия', 'рецензии', 'рецензий'))}
        </dl>
        ${ratingsTotal > 0 ? html`
          <h3 class="subheading">Распределение оценок · ${pluralCount(ratingsTotal, 'оценка', 'оценки', 'оценок')}</h3>
          <ul class="histogram">${histogram.map((b) => html`
            <li class="histogram__col">
              <span class="histogram__track"><span class="histogram__bar" style="--value: ${b.height}%"></span></span>
              <span class="histogram__label" aria-hidden="true">${b.to}</span>
              <span class="visually-hidden">${b.from}–${b.to}: ${pluralCount(b.count, 'оценка', 'оценки', 'оценок')}</span>
            </li>`)}</ul>` : ''}
      </section>

      <section class="section" aria-labelledby="book-reviews-title">
        <div class="section-heading">
          <h2 class="section-heading__title" id="book-reviews-title">Оценки и рецензии</h2>
        </div>
        ${reviews.length
          ? html`<ul class="review-list review-list--single">${reviews.map(reviewItem)}</ul>`
          : emptyState('Оценок и рецензий пока нет.')}
        ${pagination(page, hasNext)}
      </section>
    </div>`);

  const panel = $('#shelf-panel');
  if (me && entry !== undefined) {
    bindShelfPanel(panel, {
      userId: me.id,
      bookId: id,
      reload: (message) => load({ notice: message, focusPanel: true }),
    });
  }
  // после действия фокус возвращается в панель, иначе он теряется вместе с перерисованной кнопкой
  if (focusPanel) panel.focus({ preventScroll: true });
}

start({ nav: 'books', navDetail: true }, () => load());
