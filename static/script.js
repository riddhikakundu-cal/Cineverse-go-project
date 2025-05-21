function toggleModal() {
  document.getElementById('auth-modal').style.display = 'flex';
}

window.onclick = function(event) {
  if (event.target.classList.contains('modal')) {
    document.getElementById('auth-modal').style.display = 'none';
  }
}

function switchTab(tab) {
  document.getElementById('login-form').style.display = tab === 'login' ? 'block' : 'none';
  document.getElementById('signup-form').style.display = tab === 'signup' ? 'block' : 'none';
  document.getElementById('login-tab').classList.toggle('active', tab === 'login');
  document.getElementById('signup-tab').classList.toggle('active', tab === 'signup');
}

async function fetchMovies() {
  const search = document.getElementById('search').value;
  const res = await fetch(`/movies${search ? '?search=' + search : ''}`);
  const movies = await res.json();
  const list = document.getElementById('movie-list');
  list.innerHTML = movies.map(m => `
    <div class="card">
      <h3>${m.title}</h3>
      <p><strong>Genre:</strong> ${m.genre || 'N/A'}</p>
      <p><strong>Year:</strong> ${m.year || 'N/A'}</p>
      <p><strong>ISBN:</strong> ${m.isbn}</p>
      <p><strong>Director:</strong> ${m.director?.firstname || ''} ${m.director?.lastname || ''}</p>
      <div class="actions">
        <button onclick="deleteMovie('${m.id}')">Delete</button>
      </div>
    </div>
  `).join('');
}

async function deleteMovie(id) {
  if (!confirm("Are you sure you want to delete this movie?")) return;
  await fetch(`/movies/${id}`, { method: 'DELETE' });
  fetchMovies();
}

document.getElementById('movie-form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const title = document.getElementById('title').value;
  const isbn = document.getElementById('isbn').value;
  const genre = document.getElementById('genre').value;
  const year = parseInt(document.getElementById('year').value);
  const firstname = document.getElementById('firstname').value;
  const lastname = document.getElementById('lastname').value;

  await fetch('/movies', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ title, isbn, genre, year, director: { firstname, lastname } })
  });

  e.target.reset();
  fetchMovies();
});

function login() {
  const user = document.getElementById('login-username').value;
  const pass = document.getElementById('login-password').value;
  alert(`(Stub) Logged in as: ${user}`);
  document.getElementById('auth-modal').style.display = 'none';
}

function signup() {
  const user = document.getElementById('signup-username').value;
  const pass = document.getElementById('signup-password').value;
  alert(`(Stub) Signed up as: ${user}`);
  document.getElementById('auth-modal').style.display = 'none';
}

fetchMovies();
