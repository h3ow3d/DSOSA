(() => {
  // Catalogue sync button loading state
  const syncForms = document.querySelectorAll('form[action="/catalogue/sync"]');
  for (const form of syncForms) {
    form.addEventListener('submit', () => {
      const btn = form.querySelector('button');
      if (btn) { btn.textContent = 'Syncing…'; btn.disabled = true; }
    });
  }

  // Assessment scoring form — JSON POST
  const scoreForm = document.getElementById('score-form');
  const saveStatus = document.getElementById('save-status');

  function showStatus(msg, ok) {
    if (!saveStatus) return;
    saveStatus.textContent = msg;
    saveStatus.className = 'save-status ' + (ok ? 'ok' : 'err');
    setTimeout(() => { saveStatus.textContent = ''; saveStatus.className = 'save-status'; }, 4000);
  }

  if (scoreForm) {
    scoreForm.addEventListener('submit', (e) => {
      e.preventDefault();

      const scores = [];
      document.querySelectorAll('.control-row').forEach((row) => {
        const cid = row.dataset.cid;
        if (!cid) return;
        scores.push({
          control_id:    cid,
          current:       parseInt(row.querySelector('.current-select').value, 10) || 0,
          target:        parseInt(row.querySelector('.target-select').value, 10) || 0,
          not_applicable: row.querySelector('.na-check').checked,
          evidence:      row.querySelector('.evidence-area').value,
          action_notes:  row.querySelector('.action-area').value,
          priority:      row.querySelector('.priority-select').value,
          confidence:    row.querySelector('.confidence-select').value,
        });
      });

      const url = window.location.pathname + '/scores';
      fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(scores),
      })
        .then((r) => r.json())
        .then((d) => showStatus(d.message || 'Saved', true))
        .catch(() => showStatus('Error saving scores — please try again', false));
    });
  }
})();
