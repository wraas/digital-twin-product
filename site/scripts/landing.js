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

  // April Fools badge reveal
  var afBadge = document.getElementById('af-badge');
  function revealAF() {
    try { localStorage.setItem('wraas-af-revealed', '1'); } catch (e) {}
    if (afBadge) afBadge.hidden = false;
  }
  try { if (localStorage.getItem('wraas-af-revealed') === '1' && afBadge) afBadge.hidden = false; } catch (e) {}

  // Update CTA button if user has read the docs
  var docsRead = false;
  try { docsRead = localStorage.getItem('wraas-docs-read') === '1'; } catch (e) {}

  if (!docsRead) {
    // Set AF flag before navigating to rickroll
    var rickrollBtn = document.querySelector('.btn-skeptical');
    if (rickrollBtn) {
      rickrollBtn.addEventListener('click', function () { revealAF(); });
    }
  }

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
        revealAF();
        modalBody.innerHTML =
          '<div class="modal-confirmation">' +
            '<div class="modal-icon">✅</div>' +
            '<h2 id="modal-title">Access granted.</h2>' +
            '<p>No queue. No approval board. No Series A required.</p>' +

            '<div class="modal-install" role="group" aria-labelledby="install-heading">' +
              '<h3 id="install-heading" class="modal-install-heading">Install via Go</h3>' +
              '<div class="terminal modal-terminal" role="region" aria-label="Install command">' +
                '<div><span class="prompt" aria-hidden="true">$</span> <code id="install-cmd">go install github.com/wraas/digital-twin-product/cli@latest</code></div>' +
                '<button type="button" class="copy-btn" data-target="install-cmd" aria-label="Copy install command">Copy</button>' +
              '</div>' +
              '<p class="modal-alt">Or download a binary from <a href="https://github.com/wraas/digital-twin-product/releases" target="_blank" rel="noopener">GitHub Releases</a>.</p>' +
            '</div>' +

            '<div class="modal-install" role="group" aria-labelledby="verify-heading">' +
              '<h3 id="verify-heading" class="modal-install-heading">Verify</h3>' +
              '<div class="terminal modal-terminal" role="region" aria-label="Verify commands">' +
                '<div><span class="prompt" aria-hidden="true">$</span> <code>wraas init</code></div>' +
                '<div class="output">&gt; Config written to ./wraas.yml</div>' +
                '<div><span class="prompt" aria-hidden="true">$</span> <code>wraas status</code></div>' +
                '<div class="success">&gt; Engine: RUNNING</div>' +
                '<div class="warn">&gt; Latency: 113ms</div>' +
              '</div>' +
            '</div>' +

            '<hr style="border:none;border-top:1px solid rgba(0,200,255,0.15);margin:1.5rem 0">' +
            '<p class="modal-alt">Questions? The real Romain is available. His response latency is higher than 113ms.</p>' +
            '<div class="modal-contact-links">' +
              '<a href="https://github.com/rlespinasse" target="_blank" rel="noopener" class="modal-contact">GitHub &rarr;</a>' +
              '<a href="https://www.linkedin.com/in/romain-lespinasse/" target="_blank" rel="noopener" class="modal-contact">LinkedIn &rarr;</a>' +
            '</div>' +
          '</div>';

        // Wire up copy button
        var copyBtn = modalBody.querySelector('.copy-btn');
        if (copyBtn) {
          copyBtn.addEventListener('click', function () {
            var target = document.getElementById(copyBtn.getAttribute('data-target'));
            if (target && navigator.clipboard) {
              navigator.clipboard.writeText(target.textContent).then(function () {
                copyBtn.textContent = 'Copied';
                copyBtn.setAttribute('aria-label', 'Copied to clipboard');
                setTimeout(function () {
                  copyBtn.textContent = 'Copy';
                  copyBtn.setAttribute('aria-label', 'Copy install command');
                }, 2000);
              });
            }
          });
        }

        backdrop.querySelector('.modal-close').focus();
      });
    }
  }

})();

console.log('%c W.R.A.A.S. ', 'background:#00c8ff;color:#050d1a;font-weight:900;font-size:1.4rem;padding:4px 14px;border-radius:3px;letter-spacing:0.1em;');
console.log('%c We Rickroll Absolutely Anyone, Seriously. ', 'color:#00c8ff;font-size:1rem;font-weight:700;');
console.log('%c Happy April Fools 2026 ', 'color:#ffb300;font-size:0.9rem;font-style:italic;');
console.log('%c https://wraas.click ', 'color:#6a8faf;font-size:0.85rem;text-decoration:underline;');
console.log('%c You know the rules. ', 'color:#6a8faf;font-size:0.8rem;font-style:italic;');

// GoatCounter analytics
(function () {
  window.goatcounter = { path: function () { return location.pathname + location.hash; } };
  var s = document.createElement('script');
  s.async = true;
  s.src = '//gc.zgo.at/count.js';
  s.dataset.goatcounter = 'https://wraas.goatcounter.com/count';
  document.body.appendChild(s);
  document.addEventListener('click', function (e) {
    var a = e.target.closest('a[href], button');
    if (a && window.goatcounter && window.goatcounter.count) {
      var name = a.dataset.event || a.textContent.trim().substring(0, 50);
      window.goatcounter.count({ path: 'event/' + name, title: name, event: true });
    }
  });
})();
