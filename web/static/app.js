(() => {
  const forms = document.querySelectorAll('form[action="/catalogue/sync"]');
  for (const form of forms) {
    form.addEventListener("submit", () => {
      const btn = form.querySelector("button");
      if (btn) {
        btn.textContent = "Syncing...";
        btn.setAttribute("disabled", "true");
      }
    });
  }
})();
