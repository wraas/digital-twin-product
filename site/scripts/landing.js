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
