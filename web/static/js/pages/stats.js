import { api, ApiError, query } from '../api.js';
import { date, number, plural, pluralCount } from '../format.js';
import { $, cover, emptyState, html, mount, readBooks, score, setTitle, start } from '../ui.js';

const TOP = 10;
const DATE_RE = /^\d{4}-\d{2}-\d{2}$/;

const params = new URLSearchParams(location.search);
const from = params.get('from') || '';
const to = params.get('to') || '';

function stat(value, one, few, many) {
  return html`<div class="stat"><dt class="stat__label">${plural(value, one, few, many)}</dt><dd class="stat__value">${number(value)}</dd></div>`;
}

function topList(books, meta, empty) {
  if (books.length === 0) return emptyState(empty);
  return html`<ol class="top-list">${books.map((b) => html`
    <li class="top-list__item">
      ${cover(b, { small: true, extra: 'top-list__cover', decorative: true })}
      <div class="top-list__body">
        <a class="top-list__title" href="/books/${b.id}">${b.title}</a>
        <p class="top-list__meta">${meta(b)}</p>
      </div>
      ${score(b.score)}
    </li>`)}</ol>`;
}

function render({ error = '', stats = null }) {
  const period = from || to
    ? html`Период:${from && DATE_RE.test(from) ? ` с ${date(from + 'T00:00')}` : ''}${to && DATE_RE.test(to) ? ` по ${date(to + 'T00:00')}` : ''}`
    : 'За всё время';

  mount($('#app'), html`
    <div class="container page">
      <header class="page-header">
        <h1 class="page-title">Статистика</h1>
        <p class="page-subtitle">${period}</p>
      </header>

      <form class="period-form" action="/stats" method="get">
        <label class="period-form__field">С <input class="form__input" type="date" name="from" value="${from}"></label>
        <label class="period-form__field">по <input class="form__input" type="date" name="to" value="${to}"></label>
        <button class="button button--cta button--small" type="submit">Показать</button>
        ${from || to ? html`<a class="period-form__reset" href="/stats">Сбросить</a>` : ''}
      </form>

      ${error ? html`<p class="alert" role="alert">${error}</p>` : ''}

      ${stats ? html`
        <dl class="stat-grid stat-grid--large">
          ${stat(stats.users_count, 'читатель', 'читателя', 'читателей')}
          ${stat(stats.books_count, 'книга в каталоге', 'книги в каталоге', 'книг в каталоге')}
          ${stat(stats.books_on_shelves, 'книга добавлена на полки', 'книги добавлены на полки', 'книг добавлено на полки')}
          ${stat(stats.reads_count, 'книга прочитана', 'книги прочитаны', 'книг прочитано')}
        </dl>
        <div class="two-cols">
          <section class="section" aria-labelledby="top-score-title">
            <div class="section-heading">
              <h2 class="section-heading__title" id="top-score-title">Лучшие по оценкам</h2>
              <span class="section-heading__aside">за всё время</span>
            </div>
            ${topList(stats.top_by_score || [], (b) => b.author, 'Пока ни одна книга не набрала достаточно оценок.')}
          </section>
          <section class="section" aria-labelledby="top-reads-title">
            <div class="section-heading">
              <h2 class="section-heading__title" id="top-reads-title">Самые читаемые</h2>
              <span class="section-heading__aside">за всё время</span>
            </div>
            ${topList(readBooks(stats.top_by_reads, TOP).map((b) => ({ ...b, score: null })),
              (b) => html`${b.author} · ${pluralCount(b.reads_count, 'прочтение', 'прочтения', 'прочтений')}`,
              'Пока никто не отметил книгу прочитанной.')}
          </section>
        </div>` : ''}
    </div>`);
}

start({ nav: 'stats' }, async () => {
  setTitle('Статистика');

  // ошибку периода показываем прямо над формой, а не отдельной страницей
  if ((from && !DATE_RE.test(from)) || (to && !DATE_RE.test(to))) {
    render({ error: 'Даты указываются в формате ГГГГ-ММ-ДД.' });
    return;
  }

  try {
    render({ stats: await api('/stats' + query({ from, to, top: TOP })) });
  } catch (err) {
    if (!(err instanceof ApiError) || err.status !== 400) throw err;
    render({ error: 'Начало периода не может быть позже конца.' });
  }
});
