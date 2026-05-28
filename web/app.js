const API = '/api/notes';
let editingId = null;

const $ = (sel) => document.querySelector(sel);
const $$ = (sel) => document.querySelectorAll(sel);

async function apiFetch(url, opts = {}) {
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json', ...opts.headers },
    ...opts,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error || res.statusText);
  }
  return res.status !== 204 ? res.json() : null;
}

function showSnackbar(message, isError = false) {
  const el = $('#snackbar');
  el.textContent = message;
  el.className = 'snackbar show' + (isError ? ' error' : '');
  setTimeout(() => {
    el.classList.remove('show');
  }, 4000);
}

function showLoading() {
  $('#notes-list').innerHTML = '<div class="spinner-container"><div class="spinner"></div></div>';
}

async function loadNotes() {
  showLoading();
  try {
    const data = await apiFetch(API);
    renderNotes(data.notes || []);
  } catch (err) {
    $('#notes-list').innerHTML = '<div class="empty-state"><span class="material-icons empty-icon">error_outline</span><p class="empty-text">Failed to load notes</p></div>';
    showSnackbar(err.message, true);
  }
}

function renderNotes(notes) {
  const container = $('#notes-list');
  if (notes.length === 0) {
    container.innerHTML = '<div class="empty-state"><span class="material-icons empty-icon">note_add</span><p class="empty-text">No notes yet. Create one above!</p></div>';
    return;
  }
  container.innerHTML = notes.map((n, i) => `
    <div class="note" data-id="${n.id}" style="animation-delay: ${i * 0.05}s">
      <div class="note-header">
        <div class="note-title">${esc(n.title)}</div>
        <div class="note-actions">
          <button class="icon-btn" onclick="editNote('${n.id}')" title="Edit note">
            <span class="material-icons">edit</span>
          </button>
          <button class="icon-btn" onclick="deleteNote('${n.id}')" title="Delete note">
            <span class="material-icons">delete</span>
          </button>
        </div>
      </div>
      <div class="note-content">${esc(n.content)}</div>
      <span class="note-date">${new Date(n.created_at).toLocaleString()}</span>
    </div>
  `).join('');
}

function esc(str) {
  const d = document.createElement('div');
  d.textContent = str;
  return d.innerHTML;
}

function resetForm() {
  $('#note-form').reset();
  editingId = null;
  $('#form-title').textContent = 'New Note';
  $('#submit-btn').innerHTML = '<span class="material-icons" style="font-size:18px">add</span> Create';
  $('#cancel-btn').classList.add('hidden');
}

$('#cancel-btn').addEventListener('click', resetForm);

$('#note-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const title = $('#title').value.trim();
  const content = $('#content').value.trim();
  if (!title || !content) return;

  const btn = $('#submit-btn');
  const original = btn.innerHTML;
  btn.disabled = true;
  btn.innerHTML = '<div class="spinner" style="width:18px;height:18px;border-width:2px"></div>';

  try {
    if (editingId) {
      await apiFetch(`${API}/${editingId}`, {
        method: 'PUT',
        body: JSON.stringify({ title, content }),
      });
      showSnackbar('Note updated');
    } else {
      await apiFetch(API, {
        method: 'POST',
        body: JSON.stringify({ title, content }),
      });
      showSnackbar('Note created');
    }
    resetForm();
    await loadNotes();
  } catch (err) {
    showSnackbar(err.message, true);
    btn.disabled = false;
    btn.innerHTML = original;
  }
});

async function editNote(id) {
  try {
    const note = await apiFetch(`${API}/${id}`);
    $('#title').value = note.title;
    $('#content').value = note.content;
    editingId = note.id;
    $('#form-title').textContent = 'Edit Note';
    $('#submit-btn').innerHTML = '<span class="material-icons" style="font-size:18px">save</span> Update';
    $('#cancel-btn').classList.remove('hidden');
    window.scrollTo({ top: 0, behavior: 'smooth' });
  } catch (err) {
    showSnackbar(err.message, true);
  }
}

async function deleteNote(id) {
  if (!confirm('Delete this note?')) return;
  try {
    await apiFetch(`${API}/${id}`, { method: 'DELETE' });
    showSnackbar('Note deleted');
    await loadNotes();
  } catch (err) {
    showSnackbar(err.message, true);
  }
}

loadNotes();
