/* W.R.A.A.S. — We Rickroll Absolutely Anyone, Seriously. */
(function () {
  // Random pumped-up number, always ends in 0, between 10,000 and 99,990
  var n = (Math.floor(Math.random() * 8991) + 1000) * 10;
  var formatted = n.toLocaleString('en-US');
  document.getElementById('pumped-stat').innerHTML = formatted
    .split('').map(function (c) { return c === '0' ? '<span class="zero">0</span>' : c; }).join('');

  // Mobile menu toggle
  var toggle = document.querySelector('.nav-toggle');
  var menu = document.getElementById('nav-menu');
  if (toggle && menu) {
    toggle.addEventListener('click', function () {
      var expanded = toggle.getAttribute('aria-expanded') === 'true';
      toggle.setAttribute('aria-expanded', String(!expanded));
      menu.classList.toggle('is-open');
    });
    // Close menu when clicking a nav link
    menu.addEventListener('click', function (e) {
      if (e.target.tagName === 'A') {
        toggle.setAttribute('aria-expanded', 'false');
        menu.classList.remove('is-open');
      }
    });
  }

  // Update CTA button if user has read the docs
  var docsRead = false;
  try { docsRead = localStorage.getItem('wraas-docs-read') === '1'; } catch (e) {}

  if (docsRead) {
    var btn = document.querySelector('.btn-skeptical');
    if (btn) {
      btn.innerHTML = 'I <span class="gradient-text">actually</span> read the docs, so I can <span class="gradient-text">request access</span>* now';
      btn.href = '#';

      var backdrop = document.getElementById('access-modal');
      var modal = backdrop.querySelector('.modal');
      var closeBtn = backdrop.querySelector('.modal-close');
      var form = document.getElementById('access-form');
      var emailInput = document.getElementById('access-email');
      var modalBody = document.getElementById('modal-body');
      var previousFocus = null;

      function openModal() {
        previousFocus = document.activeElement;
        backdrop.classList.add('is-open');
        backdrop.setAttribute('aria-hidden', 'false');
        emailInput.focus();
      }

      function closeModal() {
        backdrop.classList.remove('is-open');
        backdrop.setAttribute('aria-hidden', 'true');
        if (previousFocus) previousFocus.focus();
      }

      btn.addEventListener('click', function (e) {
        e.preventDefault();
        openModal();
      });

      closeBtn.addEventListener('click', closeModal);

      backdrop.addEventListener('click', function (e) {
        if (e.target === backdrop) closeModal();
      });

      document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape' && backdrop.classList.contains('is-open')) {
          closeModal();
        }
        // Focus trap
        if (e.key === 'Tab' && backdrop.classList.contains('is-open')) {
          var focusable = modal.querySelectorAll('button, input, [href], [tabindex]:not([tabindex="-1"])');
          var first = focusable[0];
          var last = focusable[focusable.length - 1];
          if (e.shiftKey) {
            if (document.activeElement === first) { e.preventDefault(); last.focus(); }
          } else {
            if (document.activeElement === last) { e.preventDefault(); first.focus(); }
          }
        }
      });

      form.addEventListener('submit', function (e) {
        e.preventDefault();
        modalBody.innerHTML =
          '<div class="modal-confirmation">' +
            '<div class="modal-icon">📬</div>' +
            '<h2>Request received.</h2>' +
            '<p>Your access request has been queued. A confirmation email will be sent after our Series A funding round closes.</p>' +
            '<p>Current funding status:</p>' +
            '<span class="funding-status">Pre-seed — Estimated timeline: optimistic</span>' +
            '<p style="margin-top:1.5rem">In the meantime, WRAAS appreciates your patience. It has been noted. It has been evaluated. It has been filed.</p>' +
            '<hr style="border:none;border-top:1px solid rgba(0,200,255,0.15);margin:1.5rem 0">' +
            '<p style="font-size:0.82rem">Can\u2019t wait? You can contact the real Romain directly. Note: his availability window is shorter than the estimated funding timeline. Significantly.</p>' +
            '<div class="modal-contact-links">' +
              '<a href="https://github.com/rlespinasse" target="_blank" rel="noopener" class="modal-contact">GitHub &rarr;</a>' +
              '<a href="https://www.linkedin.com/in/romain-lespinasse/" target="_blank" rel="noopener" class="modal-contact">LinkedIn &rarr;</a>' +
            '</div>' +
          '</div>';
        backdrop.querySelector('.modal-close').focus();
      });
    }
  }

  // Shuffle capabilities cards
  var grid = document.querySelector('.features-grid');
  var cards = [].slice.call(grid.children);
  for (var i = cards.length - 1; i > 0; i--) {
    var j = Math.floor(Math.random() * (i + 1));
    grid.appendChild(cards[j]);
    cards.splice(j, 1);
  }
})();

console.log('%c W.R.A.A.S. ', 'background:#00c8ff;color:#050d1a;font-weight:900;font-size:1.4rem;padding:4px 14px;border-radius:3px;letter-spacing:0.1em;');
console.log('%c We Rickroll Absolutely Anyone, Seriously. ', 'color:#00c8ff;font-size:1rem;font-weight:700;');
console.log('%c Happy April Fools 2026 ', 'color:#ffb300;font-size:0.9rem;font-style:italic;');
console.log('%c https://wraas.click ', 'color:#6a8faf;font-size:0.85rem;text-decoration:underline;');
console.log('%c You know the rules. ', 'color:#6a8faf;font-size:0.8rem;font-style:italic;');
