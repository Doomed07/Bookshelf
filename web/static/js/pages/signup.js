import { api, ApiError } from '../api.js';
import { $, renderLayout, setTitle } from '../ui.js';

// Те же правила, что в домене (internal/core/domain/user.go): проверяем заранее,
// чтобы показать ошибку у каждого поля, а не только первую от сервера.
const USERNAME_RE = /^[A-Za-z0-9_]{3,30}$/;
const EMAIL_RE = /^[a-z0-9._+-]+@([a-z0-9-]+\.)+[a-z]{2,}$/;

const USERNAME_HINT = 'От 3 до 30 символов: латинские буквы, цифры и _.';
const USERNAME_ERROR = USERNAME_HINT;
const EMAIL_ERROR = 'Введите email в формате name@example.com.';
const CONFLICT_ERROR = 'Имя пользователя или email уже заняты.';
const INVALID_ERROR = 'Проверьте имя пользователя и email.';
const UNKNOWN_ERROR = 'Не удалось создать профиль. Попробуйте ещё раз чуть позже.';

renderLayout({ nav: 'signup' });
setTitle('Регистрация');

const form = $('#signup-form');
const alertBox = $('#signup-alert');
const username = $('#signup-username');
const email = $('#signup-email');
const usernameHint = $('#signup-username-hint');
const emailHint = $('#signup-email-hint');

function setFieldError(input, hint, message, defaultText = '') {
  input.toggleAttribute('aria-invalid', Boolean(message));
  hint.textContent = message || defaultText;
  hint.className = message ? 'form__error' : 'form__hint';
  hint.hidden = !message && !defaultText;
}

function showAlert(message) {
  alertBox.textContent = message;
  alertBox.hidden = !message;
}

form.addEventListener('submit', async (event) => {
  event.preventDefault();
  showAlert('');

  const name = username.value.trim();
  const mail = email.value.trim().toLowerCase();
  username.value = name;
  email.value = mail;

  const usernameOK = USERNAME_RE.test(name);
  const emailOK = mail.length <= 254 && EMAIL_RE.test(mail);
  setFieldError(username, usernameHint, usernameOK ? '' : USERNAME_ERROR, USERNAME_HINT);
  setFieldError(email, emailHint, emailOK ? '' : EMAIL_ERROR);
  if (!usernameOK || !emailOK) {
    (usernameOK ? email : username).focus();
    return;
  }

  const button = form.querySelector('button[type="submit"]');
  button.disabled = true;
  try {
    const user = await api('/users', { method: 'POST', body: { username: name, email: mail } });
    location.href = `/users/${user.id}`;
  } catch (err) {
    const status = err instanceof ApiError ? err.status : 0;
    showAlert(status === 409 ? CONFLICT_ERROR : status === 400 ? INVALID_ERROR : UNKNOWN_ERROR);
    button.disabled = false;
  }
});
