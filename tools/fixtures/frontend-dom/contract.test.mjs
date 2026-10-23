import { test, afterEach, after } from 'node:test';
import assert from 'node:assert/strict';
import { JSDOM } from 'jsdom';

// Set the DOM before React DOM/Testing Library imports bind their event environment.
const dom = new JSDOM('<!doctype html><html><body></body></html>', { url: 'https://trial.example.invalid' });
const previous = new Map();
for (const [name, value] of Object.entries({ window: dom.window, document: dom.window.document,
  navigator: dom.window.navigator, HTMLElement: dom.window.HTMLElement,
  IS_REACT_ACT_ENVIRONMENT: true })) {
  previous.set(name, Object.getOwnPropertyDescriptor(globalThis, name));
  Object.defineProperty(globalThis, name, { configurable: true, writable: true, value });
}
const { createElement: h } = await import('react');
const { render, screen, fireEvent, act, cleanup, within } = await import('@testing-library/react');
const { NotesPanel } = await import('./NotesPanel.mjs');
const memory = dom.window.localStorage;
memory.setItem('untouched', 'synthetic-only');
const storageSnapshot = () => Array.from({ length: memory.length }, (_, i) => memory.key(i))
  .sort().map(key => [key, memory.getItem(key)]);
const initialStorage = storageSnapshot();

function deferred() {
  let resolve, reject;
  const promise = new Promise((yes, no) => { resolve = yes; reject = no; });
  return { promise, resolve, reject };
}
async function settle(pending, value, reject = false) {
  await act(async () => { pending[reject ? 'reject' : 'resolve'](value); });
}
function listAdapter(result) {
  return { list: async () => result, create: async () => { throw new Error('Unexpected create'); } };
}
function noConflictingListStates() {
  assert.equal(screen.queryByText('Loading notes...'), null);
  assert.equal(screen.queryByText('Unable to load notes. Retry.'), null);
}

afterEach(() => {
  try { assert.deepEqual(storageSnapshot(), initialStorage, 'Stored browser data changed'); }
  finally { cleanup(); }
  assert.equal(document.body.childElementCount, 0, 'Rendered roots leaked');
});
after(() => {
  dom.window.close();
  for (const [name, descriptor] of previous) {
    if (descriptor) Object.defineProperty(globalThis, name, descriptor);
    else delete globalThis[name];
  }
});

test('loading has a visible status, no premature result and no create effect', async () => {
  const request = deferred();
  let loads = 0, creates = 0;
  render(h(NotesPanel, { client: { list: () => { loads++; return request.promise; },
    create: async () => { creates++; } } }));
  assert.equal(screen.getByRole('status').textContent, 'Loading notes...');
  assert.equal(screen.queryByRole('list', { name: 'Notes' }), null);
  assert.equal(screen.queryByText('No notes yet.'), null);
  assert.equal(screen.queryByRole('alert'), null);
  assert.equal(loads, 1);
  assert.equal(creates, 0);
  await settle(request, []);
});

test('list failure is safe and retry clears the error before a visible success', async () => {
  const first = deferred(), second = deferred();
  let calls = 0;
  render(h(NotesPanel, { client: { list: () => (++calls === 1 ? first : second).promise,
    create: async () => { throw new Error('Unexpected create'); } } }));
  assert.equal(screen.getByRole('status').textContent, 'Loading notes...');
  await settle(first, new Error('synthetic-private-diagnostic'), true);
  assert.equal(screen.getByRole('alert').textContent, 'Unable to load notes. Retry.');
  assert.equal(screen.queryByText('Loading notes...'), null);
  assert.equal(screen.queryByText('No notes yet.'), null);
  assert.equal(screen.queryByRole('list', { name: 'Notes' }), null);
  assert.ok(!document.body.textContent.includes('synthetic-private-diagnostic'));
  fireEvent.click(screen.getByRole('button', { name: 'Retry loading' }));
  assert.equal(calls, 2);
  assert.equal(screen.queryByRole('alert'), null);
  assert.equal(screen.getByRole('status').textContent, 'Loading notes...');
  await settle(second, [{ id: 'a', text: 'Recovered note' }]);
  assert.equal(screen.getByRole('list', { name: 'Notes' }).textContent, 'Recovered note');
  assert.equal(screen.queryByRole('button', { name: 'Retry loading' }), null);
  assert.equal(screen.queryByRole('alert'), null);
  noConflictingListStates();
});

test('empty response is distinct from loading, failure and a populated list', async () => {
  render(h(NotesPanel, { client: listAdapter([]) }));
  assert.equal((await screen.findByText('No notes yet.')).textContent, 'No notes yet.');
  assert.equal(screen.queryByRole('list', { name: 'Notes' }), null);
  assert.equal(screen.queryByRole('alert'), null);
  noConflictingListStates();
});

test('success renders the requested records as text, without stale states', async () => {
  const literal = '<img src=x onerror=synthetic>';
  render(h(NotesPanel, { client: listAdapter([{ id: 'a', text: 'First note' }, { id: 'b', text: literal }]) }));
  const list = await screen.findByRole('list', { name: 'Notes' });
  assert.deepEqual(within(list).getAllByRole('listitem').map(row => row.textContent), ['First note', literal]);
  assert.equal(screen.queryByRole('img'), null);
  assert.equal(screen.queryByText('No notes yet.'), null);
  assert.equal(screen.queryByRole('alert'), null);
  noConflictingListStates();
});

test('invalid form input shows a field error without API writes or stored changes', async () => {
  const records = [];
  let creates = 0;
  const client = { list: async () => structuredClone(records), create: async (payload) => {
    creates++; records.push({ id: 'a', text: payload.text }); return records.at(-1);
  } };
  render(h(NotesPanel, { client }));
  await screen.findByText('No notes yet.');
  const before = await client.list();
  const input = screen.getByRole('textbox', { name: 'Note' });
  fireEvent.change(input, { target: { value: '   ' } });
  fireEvent.click(screen.getByRole('button', { name: 'Add note' }));
  assert.equal(screen.getByRole('alert').textContent, 'Enter a note.');
  assert.equal(creates, 0);
  assert.deepEqual(await client.list(), before);
  assert.equal(input.value, '   ');
  assert.equal(input.disabled, false);
  assert.equal(screen.getByRole('button', { name: 'Add note' }).disabled, false);
  assert.equal(screen.queryByText('Saving note...'), null);
  assert.equal(screen.queryByText('Note saved.'), null);
});

test('valid form gates duplicates, retains failed input and recovers with one stored result', async () => {
  const first = deferred(), second = deferred();
  const records = [], payloads = [];
  const client = { list: async () => structuredClone(records), create: async (payload) => {
    payloads.push(structuredClone(payload));
    const result = await (payloads.length === 1 ? first : second).promise;
    records.push(result); return result;
  } };
  render(h(NotesPanel, { client }));
  await screen.findByText('No notes yet.');
  const input = screen.getByRole('textbox', { name: 'Note' });
  const button = screen.getByRole('button', { name: 'Add note' });
  fireEvent.change(input, { target: { value: '  New note  ' } });
  fireEvent.click(button);
  assert.deepEqual(payloads, [{ text: 'New note' }]);
  assert.equal(input.disabled, true);
  assert.equal(button.disabled, true);
  assert.equal(screen.getByText('Saving note...').getAttribute('role'), 'status');
  assert.equal(screen.queryByText('Note saved.'), null);
  fireEvent.click(button);
  assert.equal(payloads.length, 1);
  assert.deepEqual(await client.list(), []);
  await settle(first, new Error('synthetic-secret-detail'), true);
  assert.equal(screen.getByRole('alert').textContent, 'Unable to save note. Try again.');
  assert.ok(!document.body.textContent.includes('synthetic-secret-detail'));
  assert.equal(input.value, '  New note  ');
  assert.equal(input.disabled, false);
  assert.equal(button.disabled, false);
  assert.equal(screen.queryByText('Saving note...'), null);
  assert.equal(screen.queryByText('Note saved.'), null);
  assert.deepEqual(await client.list(), []);
  fireEvent.click(button);
  assert.deepEqual(payloads, [{ text: 'New note' }, { text: 'New note' }]);
  assert.equal(screen.queryByRole('alert'), null);
  assert.equal(button.disabled, true);
  await settle(second, { id: 'created', text: 'New note' });
  assert.equal(screen.getByRole('status').textContent, 'Note saved.');
  assert.equal(screen.getByRole('list', { name: 'Notes' }).textContent, 'New note');
  assert.equal(screen.queryByText('No notes yet.'), null);
  assert.equal(screen.queryByRole('alert'), null);
  assert.equal(screen.queryByText('Saving note...'), null);
  assert.equal(input.value, '');
  assert.equal(input.disabled, false);
  assert.equal(button.disabled, false);
  assert.deepEqual(await client.list(), [{ id: 'created', text: 'New note' }]);
});
