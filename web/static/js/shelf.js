// Своя полка: запросы к API, тексты ошибок, панель на странице книги и кнопки под обложками в профиле.
//
// Все действия идут от имени вошедшего (cookie), а сервер сам не даст менять чужую полку (403).
// Пользовательский текст (рецензия) подставляется только через html`…`, то есть экранируется.

import { api, ApiError } from './api.js';
import { html, mount } from './ui.js';
import { date, isoDate } from './format.js';

export const MAX_REVIEW = 5000;

// ---------- Запросы ----------

const shelfPath = (userId) => `/users/${userId}/bookshelf`;

// getShelfBook — запись о книге на полке пользователя или null, если книги на полке нет.
export async function getShelfBook(userId, bookId) {
  try {
    return await api(`${shelfPath(userId)}/${bookId}`);
  } catch (err) {
    if (err instanceof ApiError && err.status === 404) return null;
    throw err;
  }
}

export const addToShelf = (userId, bookId) => api(shelfPath(userId), { method: 'POST', body: { book_id: bookId } });
export const patchShelfBook = (userId, bookId, patch) => api(`${shelfPath(userId)}/${bookId}`, { method: 'PATCH', body: patch });
export const removeFromShelf = (userId, bookId) => api(`${shelfPath(userId)}/${bookId}`, { method: 'DELETE' });

// ---------- Тексты ошибок ----------

// actionMessage — что показать человеку, если действие с полкой не удалось.
// kind === 'add' отличает «книга уже на полке» от «оценку нельзя поставить непрочитанной».
export function actionMessage(err, kind = 'patch') {
  const status = err instanceof ApiError ? err.status : 0;
  switch (status) {
    case 401: {
      const next = encodeURIComponent(location.pathname + location.search);
      return html`Сессия закончилась. <a href="/login?next=${next}">Войдите снова</a>.`;
    }
    case 403: return 'Это можно делать только на своей полке.';
    case 404: return 'Этой книги нет на полке. Обновите страницу.';
    case 409: return kind === 'add'
      ? 'Эта книга уже на вашей полке. Обновите страницу.'
      : 'Оценку и рецензию можно поставить только прочитанной книге.';
    case 400: return `Проверьте оценку (от 1 до 100) и рецензию (до ${MAX_REVIEW} символов).`;
    default: return 'Не удалось выполнить действие. Попробуйте ещё раз.';
  }
}

// ---------- Разметка ----------

// confirmButton — кнопка с подтверждением: первый клик меняет её на «Точно? Да / Отмена».
function confirmButton({ action, label, text }) {
  return html`<span class="confirm" data-confirm><button class="button button--ghost button--small" type="button" data-action="${action}" data-confirm-text="${text}">${label}</button></span>`;
}

const STATUS_AREA = html`<p class="shelf-message" data-shelf-status role="status" aria-live="polite"></p>`;

function ratingField(rating) {
  const empty = rating === null || rating === undefined;
  return html`
    <div class="form__field">
      <label class="form__label" for="rating-number">Моя оценка</label>
      <div class="rating-field">
        <input class="rating-field__range" id="rating-range" type="range" min="1" max="100" value="${empty ? 50 : rating}" aria-label="Оценка ползунком, от 1 до 100"${empty ? html` data-empty` : ''}>
        <input class="form__input rating-field__number" id="rating-number" type="number" min="1" max="100" step="1" inputmode="numeric" placeholder="—" value="${empty ? '' : rating}">
        <span class="rating-field__max" aria-hidden="true">/ 100</span>
        <button class="button button--ghost button--small" type="button" data-action="clear-rating">Убрать оценку</button>
      </div>
      <p class="form__hint">Чтобы оставить книгу без оценки, очистите поле и нажмите «Сохранить».</p>
    </div>`;
}

function reviewField(review) {
  const text = review || '';
  return html`
    <div class="form__field">
      <label class="form__label" for="review-text">Рецензия</label>
      <textarea class="form__input form__textarea" id="review-text" rows="5" maxlength="${MAX_REVIEW}" placeholder="Что вы думаете об этой книге?">${text}</textarea>
      <p class="form__hint"><span id="review-count">${text.length}</span> из ${MAX_REVIEW}</p>
    </div>`;
}

function guestPanel() {
  const next = encodeURIComponent(location.pathname);
  return html`
    <p class="shelf-panel__text">Войдите, чтобы добавить книгу на свою полку, поставить оценку и написать рецензию.</p>
    <p class="shelf-actions">
      <a class="button button--cta button--small" href="/login?next=${next}">Войти</a>
      <a class="button button--ghost button--small" href="/signup?next=${next}">Зарегистрироваться</a>
    </p>`;
}

function newPanel() {
  return html`
    <p class="shelf-panel__text">Этой книги нет на вашей полке.</p>
    <p class="shelf-actions">
      <button class="button button--cta button--small" type="button" data-action="add">Хочу прочитать</button>
      <button class="button button--ghost button--small" type="button" data-action="add-read">Уже прочитал</button>
    </p>`;
}

function wantPanel() {
  return html`
    <p class="shelf-status"><span class="shelf-badge">Хочу прочитать</span></p>
    <p class="shelf-actions">
      <button class="button button--cta button--small" type="button" data-action="mark-read">Отметить прочитанной</button>
      ${confirmButton({ action: 'remove', label: 'Убрать с полки', text: 'Книга уйдёт с вашей полки.' })}
    </p>`;
}

function readPanel(shelf) {
  return html`
    <p class="shelf-status">
      <span class="shelf-badge shelf-badge--read">Прочитано</span>
      ${shelf.read_at ? html`<time datetime="${isoDate(shelf.read_at)}">${date(shelf.read_at)}</time>` : ''}
    </p>
    <form class="shelf-form" data-shelf-form novalidate>
      ${ratingField(shelf.rating)}
      ${reviewField(shelf.review)}
      <p class="shelf-actions"><button class="button button--cta button--small" type="submit">Сохранить</button></p>
    </form>
    <p class="shelf-actions shelf-actions--secondary">
      ${confirmButton({ action: 'unread', label: 'Вернуть в «хочу прочитать»', text: 'Оценка и рецензия будут удалены.' })}
      ${confirmButton({ action: 'remove', label: 'Убрать с полки', text: 'Книга, оценка и рецензия уйдут с вашей полки.' })}
    </p>`;
}

// shelfPanel — панель «Моя полка» на странице книги.
// me — вошедший или null; entry — запись о книге на полке, null (нет на полке) или undefined (не удалось загрузить).
export function shelfPanel({ me, entry, notice = '' }) {
  let body;
  if (!me) body = guestPanel();
  else if (entry === undefined) body = html`<p class="shelf-panel__text">Не удалось загрузить вашу полку. Обновите страницу.</p>`;
  else if (entry === null) body = newPanel();
  else if (entry.shelf.read) body = readPanel(entry.shelf);
  else body = wantPanel();

  return html`
    <section class="shelf-panel" aria-labelledby="shelf-panel-title">
      <h2 class="shelf-panel__title" id="shelf-panel-title">Моя полка</h2>
      ${body}
      ${me ? STATUS_AREA : ''}
      ${notice ? html`<p class="shelf-message shelf-message--error" role="alert">${notice}</p>` : ''}
    </section>`;
}

// cardActions — кнопки под обложкой на своей полке в профиле.
export function cardActions(shelf, book) {
  return html`
    <div class="shelf-actions shelf-actions--compact" data-book-id="${book.id}">
      ${shelf.read ? '' : html`<button class="button button--ghost button--small" type="button" data-action="mark-read">Прочитано</button>`}
      ${confirmButton({ action: 'remove', label: 'Убрать', text: 'Убрать с полки?' })}
    </div>`;
}

// ---------- Поведение ----------

function say(root, message, isError = false) {
  const area = root.querySelector('[data-shelf-status]');
  if (!area) return;
  mount(area, message);
  area.classList.toggle('shelf-message--error', isError && Boolean(message));
}

// Подтверждение: «Убрать» → «Точно? [Да] [Отмена]». Возвращает true, если клик обработан здесь.
function handleConfirm(button) {
  const wrap = button.closest('[data-confirm]');
  if (!wrap) return false;

  if (button.dataset.action === 'cancel') {
    wrap.replaceChildren(...wrap._original);
    wrap.querySelector('button')?.focus();
    return true;
  }

  if (button.dataset.confirmText && !button.dataset.confirmed) {
    wrap._original = [...wrap.childNodes];
    const text = document.createElement('span');
    text.className = 'confirm__text';
    text.textContent = button.dataset.confirmText;

    const yes = button.cloneNode(false);
    yes.dataset.confirmed = '1';
    yes.className = 'button button--danger button--small';
    yes.textContent = 'Да';

    const cancel = document.createElement('button');
    cancel.type = 'button';
    cancel.className = 'button button--ghost button--small';
    cancel.dataset.action = 'cancel';
    cancel.textContent = 'Отмена';

    wrap.replaceChildren(text, yes, cancel);
    cancel.focus(); // безопасный выбор по умолчанию
    return true;
  }
  return false;
}

// makeRunner выполняет действие: блокирует кнопки на время запроса (защита от двойного клика),
// затем перерисовывает страницу. Ошибку показывает рядом с кнопками.
function makeRunner(root, reload) {
  let busy = false;
  const controls = () => root.querySelectorAll('button, input, textarea');

  return async function run(task, kind) {
    if (busy) return;
    busy = true;
    controls().forEach((el) => { el.disabled = true; });
    say(root, '');

    let failure = null;
    try {
      // task возвращает ответ API, а ошибкой считается только Error — её возвращает
      // add-read, если книга добавилась, но отметить её прочитанной не удалось
      const result = await task();
      failure = result instanceof Error ? result : null;
    } catch (err) {
      busy = false;
      controls().forEach((el) => { el.disabled = false; });
      say(root, actionMessage(err, kind), true);
      return;
    }

    try {
      await reload(failure ? actionMessage(failure, 'patch') : '');
    } catch {
      busy = false;
      controls().forEach((el) => { el.disabled = false; });
      say(root, 'Изменение сохранено, но страницу не удалось обновить. Обновите её вручную.', true);
    }
  };
}

function parseRating(raw) {
  const text = raw.trim();
  if (text === '') return { value: null };
  const n = Number(text);
  if (!Number.isInteger(n) || n < 1 || n > 100) return { error: 'Оценка — целое число от 1 до 100.' };
  return { value: n };
}

// bindShelfPanel оживляет панель на странице книги.
export function bindShelfPanel(root, { userId, bookId, reload }) {
  const run = makeRunner(root, reload);
  const range = root.querySelector('#rating-range');
  const number = root.querySelector('#rating-number');
  const review = root.querySelector('#review-text');
  const counter = root.querySelector('#review-count');

  // ползунок и число показывают одно и то же; пустое число — «без оценки»
  range?.addEventListener('input', () => {
    number.value = range.value;
    range.removeAttribute('data-empty');
  });
  number?.addEventListener('input', () => {
    const n = Number(number.value);
    if (number.value.trim() === '') range.setAttribute('data-empty', '');
    else if (Number.isFinite(n)) {
      range.value = Math.min(100, Math.max(1, Math.round(n)));
      range.removeAttribute('data-empty');
    }
  });
  review?.addEventListener('input', () => {
    counter.textContent = review.value.length;
  });

  root.addEventListener('click', (event) => {
    const button = event.target.closest('button[data-action]');
    if (!button || !root.contains(button) || button.disabled) return;
    if (handleConfirm(button)) return;

    switch (button.dataset.action) {
      case 'add':
        run(() => addToShelf(userId, bookId), 'add');
        break;
      case 'add-read':
        run(async () => {
          await addToShelf(userId, bookId);
          // если вторая половина не удалась, книга остаётся в «хочу прочитать», а причину покажем после обновления
          try {
            await patchShelfBook(userId, bookId, { read: true });
          } catch (err) {
            return err;
          }
          return null;
        }, 'add');
        break;
      case 'mark-read':
        run(() => patchShelfBook(userId, bookId, { read: true }));
        break;
      case 'unread':
        run(() => patchShelfBook(userId, bookId, { read: false }));
        break;
      case 'remove':
        run(() => removeFromShelf(userId, bookId));
        break;
      case 'clear-rating':
        number.value = '';
        range.setAttribute('data-empty', '');
        number.focus();
        break;
      default:
    }
  });

  root.querySelector('[data-shelf-form]')?.addEventListener('submit', (event) => {
    event.preventDefault();
    const rating = parseRating(number.value);
    if (rating.error) {
      say(root, rating.error, true);
      number.setAttribute('aria-invalid', 'true');
      number.focus();
      return;
    }
    number.removeAttribute('aria-invalid');
    const text = review.value.trim();
    run(() => patchShelfBook(userId, bookId, { rating: rating.value, review: text === '' ? null : text }));
  });
}

// bindCardActions оживляет кнопки под обложками на своей полке в профиле.
export function bindCardActions(root, { userId, reload }) {
  const run = makeRunner(root, reload);

  root.addEventListener('click', (event) => {
    const button = event.target.closest('button[data-action]');
    if (!button || !root.contains(button) || button.disabled) return;
    if (handleConfirm(button)) return;

    const bookId = Number(button.closest('[data-book-id]')?.dataset.bookId);
    if (!bookId) return;
    if (button.dataset.action === 'mark-read') run(() => patchShelfBook(userId, bookId, { read: true }));
    if (button.dataset.action === 'remove') run(() => removeFromShelf(userId, bookId));
  });
}
