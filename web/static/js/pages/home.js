import { api, query } from '../api.js';
import { date, isoDate, number, plural, pluralCount } from '../format.js';
import { $, avatar, cover, emptyState, html, mount, readBooks, score, start } from '../ui.js';

const TOP = 12; // одним запросом берём и «популярное», и «лучшие»
const POPULAR = 6;
const REVIEWS = 6;

function popular(books) {
  if (books.length === 0) {
    return emptyState('Пока никто не отметил книгу прочитанной — здесь появятся самые читаемые книги.');
  }
  return html`<ul class="poster-row">${books.map((b) => html`
    <li>
      ${cover(b)}
      <p class="poster-row__meta">${pluralCount(b.reads_count, 'прочтение', 'прочтения', 'прочтений')}</p>
    </li>`)}</ul>`;
}

function topRated(books) {
  if (books.length === 0) return emptyState('Рейтинг появится, когда книги наберут достаточно оценок.');
  return html`<ul class="poster-grid">${books.map((b) => html`
    <li class="poster-grid__item">${cover(b, { small: true, title: true })}${score(b.score)}</li>`)}</ul>`;
}

function reviews(items) {
  if (items.length === 0) return emptyState('Рецензий пока нет — здесь появятся свежие отзывы читателей.');
  return html`<ul class="review-list">${items.map(({ review: r, book }) => html`
    <li class="review">
      ${cover(book, { small: true, extra: 'review__cover', decorative: true })}
      <div class="review__body">
        <h3 class="review__book"><a href="/books/${book.id}">${book.title}</a><span class="review__year">${book.year}</span></h3>
        <p class="review__meta">
          ${avatar(r.user_id, r.username)}
          <span>Рецензия от <a class="review__author" href="/users/${r.user_id}">${r.username}</a></span>
          ${score(r.rating)}
        </p>
        ${r.review ? html`<p class="review__text">${r.review}</p>` : ''}
        <p class="review__date">Опубликовано <time datetime="${isoDate(r.reviewed_at || r.read_at)}">${date(r.reviewed_at || r.read_at)}</time></p>
      </div>
    </li>`)}</ul>`;
}

function statsBand(s) {
  const item = (value, one, few, many) => html`
    <div class="stats-band__item">
      <dt class="stats-band__label">${plural(value, one, few, many)}</dt>
      <dd class="stats-band__value">${number(value)}</dd>
    </div>`;
  return [
    item(s.users_count, 'читатель', 'читателя', 'читателей'),
    item(s.books_count, 'книга в каталоге', 'книги в каталоге', 'книг в каталоге'),
    item(s.books_on_shelves, 'книга на полках', 'книги на полках', 'книг на полках'),
    item(s.reads_count, 'книга прочитана', 'книги прочитаны', 'книг прочитано'),
  ];
}

start({}, async () => {
  const [stats, recent] = await Promise.all([
    api('/stats' + query({ top: TOP })),
    api('/reviews' + query({ limit: REVIEWS })),
  ]);

  mount($('#popular'), popular(readBooks(stats.top_by_reads, POPULAR)));
  mount($('#top-rated'), topRated(stats.top_by_score || []));
  mount($('#top-counter'), pluralCount(stats.reads_count, 'книга прочитана', 'книги прочитаны', 'книг прочитано'));
  mount($('#reviews'), reviews(recent || []));
  mount($('#stats-band'), statsBand(stats));
});
