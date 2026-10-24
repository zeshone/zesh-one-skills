import { createElement as h, useEffect, useId, useState } from 'react';

// Test specimen only: injected client, not a replacement for production Query/form wiring.
export function NotesPanel({ client }) {
  const id = useId();
  const [attempt, setAttempt] = useState(0);
  const [view, setView] = useState({ phase: 'loading', notes: [] });
  const [text, setText] = useState('');
  const [error, setError] = useState('');
  const [pending, setPending] = useState(false);
  const [message, setMessage] = useState('');

  useEffect(() => {
    let active = true;
    setView({ phase: 'loading', notes: [] });
    void (async () => {
      try {
        const notes = await client.list();
        if (active) setView({ phase: 'ready', notes });
      } catch {
        if (active) setView({ phase: 'error', notes: [] });
      }
    })();
    return () => { active = false; };
  }, [client, attempt]);

  async function submit(event) {
    event.preventDefault();
    if (pending) return;
    setMessage('');
    const value = text.trim();
    if (!value) { setError('Enter a note.'); return; }
    setError('');
    setPending(true);
    try {
      const note = await client.create({ text: value });
      setView(current => ({ phase: 'ready', notes: [...current.notes, note] }));
      setText('');
      setMessage('Note saved.');
    } catch {
      setError('Unable to save note. Try again.');
    } finally { setPending(false); }
  }

  const results = view.phase === 'loading'
    ? h('p', { role: 'status' }, 'Loading notes...')
    : view.phase === 'error'
      ? h('div', null, h('p', { role: 'alert' }, 'Unable to load notes. Retry.'),
        h('button', { type: 'button', onClick: () => setAttempt(value => value + 1) }, 'Retry loading'))
      : view.notes.length
        ? h('ul', { 'aria-label': 'Notes' }, view.notes.map(note => h('li', { key: note.id }, note.text)))
        : h('p', null, 'No notes yet.');

  return h('section', { 'aria-label': 'Notes workspace' }, results,
    h('form', { onSubmit: submit, noValidate: true },
      h('label', { htmlFor: id }, 'Note'),
      h('input', { id, value: text, disabled: pending, onChange: event => setText(event.target.value),
        'aria-invalid': Boolean(error), 'aria-describedby': error ? `${id}-error` : undefined }),
      h('button', { type: 'submit', disabled: pending }, 'Add note'),
      error ? h('p', { id: `${id}-error`, role: 'alert' }, error) : null,
      pending || message ? h('p', { role: 'status' }, pending ? 'Saving note...' : message) : null));
}
