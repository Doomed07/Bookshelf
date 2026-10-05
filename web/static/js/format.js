// Форматирование для интерфейса: склонения, числа, даты, классы оценок и обложек.

export const NBSP = ' ';

const MONTHS_GENITIVE = [
  'января', 'февраля', 'марта', 'апреля', 'мая', 'июня',
  'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря',
];

// пороги цвета оценки 1–100: золотая, индиго, серая
const SCORE_HIGH = 80;
const SCORE_MID = 60;

// plural выбирает форму слова под число: 1 книга, 2 книги, 5 книг, 11 книг, 21 книга.
export function plural(n, one, few, many) {
  n = Math.abs(n);
  if (n % 100 >= 11 && n % 100 <= 14) return many;
  switch (n % 10) {
    case 1: return one;
    case 2: case 3: case 4: return few;
    default: return many;
  }
}

// pluralCount — число вместе со словом: «1 022 книги».
export function pluralCount(n, one, few, many) {
  return number(n) + NBSP + plural(n, one, few, many);
}

// number делит число на разряды неразрывным пробелом: 1022 → «1 022».
export function number(n) {
  const sign = n < 0 ? '-' : '';
  const digits = String(Math.abs(Math.trunc(n)));
  let head = digits.length % 3 || 3;
  let out = digits.slice(0, head);
  for (let i = head; i < digits.length; i += 3) out += NBSP + digits.slice(i, i + 3);
  return sign + out;
}

// date — дата по-русски в часовом поясе браузера: «2 октября 2026».
export function date(value) {
  const d = new Date(value);
  return `${d.getDate()}${NBSP}${MONTHS_GENITIVE[d.getMonth()]} ${d.getFullYear()}`;
}

// isoDate — дата для атрибута <time datetime>.
export function isoDate(value) {
  const d = new Date(value);
  const pad = (x) => String(x).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}

// decimal — один знак после запятой: 87.5 → «87,5», 80 → «80».
export function decimal(v) {
  if (v == null) return '—';
  return v.toFixed(1).replace(/\.0$/, '').replace('.', ',');
}

export function percent(v) {
  if (v == null) return '—';
  return Math.round(v) + NBSP + '%';
}

// duration — время чтения по числу часов: «меньше часа», «5 часов», «12 дней».
export function duration(hours) {
  if (hours == null) return '—';
  if (hours < 1) return 'меньше часа';
  if (hours < 48) return pluralCount(Math.round(hours), 'час', 'часа', 'часов');
  return pluralCount(Math.round(hours / 24), 'день', 'дня', 'дней');
}

export function scoreClass(score) {
  if (score >= SCORE_HIGH) return 'score--high';
  if (score >= SCORE_MID) return 'score--mid';
  return 'score--low';
}

const mod = (a, b) => ((a % b) + b) % b;

// coverClass — постоянный цвет обложки книги, см. .cover--0…7 в main.css.
export function coverClass(bookId) {
  return 'cover--' + mod(bookId, 8);
}

// avatarClass — постоянный цвет аватара читателя, см. .avatar--0…5 в main.css.
export function avatarClass(userId) {
  return 'avatar--' + mod(userId, 6);
}

export function initial(name) {
  const first = [...(name || '')][0];
  return first ? first.toUpperCase() : '?';
}
