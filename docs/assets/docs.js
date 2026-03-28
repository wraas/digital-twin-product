/* W.R.A.A.S. — We Rickroll Absolutely Anyone, Seriously. */
console.log('%c W.R.A.A.S. ', 'background:#00c8ff;color:#050d1a;font-weight:900;font-size:1.4rem;padding:4px 14px;border-radius:3px;letter-spacing:0.1em;');
console.log('%c We Rickroll Absolutely Anyone, Seriously. ', 'color:#0065b3;font-size:1rem;font-weight:700;');
console.log('%c Happy April Fools 2026 ', 'color:#c47d00;font-size:0.9rem;font-style:italic;');
console.log('%c https://wraas.click ', 'color:#666;font-size:0.85rem;text-decoration:underline;');
console.log('%c You know the rules. ', 'color:#999;font-size:0.8rem;font-style:italic;');

// Mobile nav toggle
var navToggle = document.querySelector('.nav-toggle');
var navPanel = document.querySelector('.nav-panel');
if (navToggle && navPanel) {
  navToggle.addEventListener('click', function () {
    var expanded = navToggle.getAttribute('aria-expanded') === 'true';
    navToggle.setAttribute('aria-expanded', String(!expanded));
    navPanel.classList.toggle('is-open');
  });
  // Close nav when clicking a link
  navPanel.addEventListener('click', function (e) {
    if (e.target.tagName === 'A') {
      navToggle.setAttribute('aria-expanded', 'false');
      navPanel.classList.remove('is-open');
    }
  });
}

// TOC scroll highlight
var tocLinks = document.querySelectorAll('.toc-link');
var sections = [].slice.call(tocLinks).map(function (l) {
  return document.querySelector(l.getAttribute('href'));
}).filter(Boolean);

var observer = new IntersectionObserver(function (entries) {
  entries.forEach(function (entry) {
    if (entry.isIntersecting) {
      tocLinks.forEach(function (l) { l.classList.remove('is-active'); });
      var active = [].slice.call(tocLinks).find(function (l) {
        return l.getAttribute('href') === '#' + entry.target.id;
      });
      if (active) active.classList.add('is-active');
    }
  });
}, { rootMargin: '-20% 0px -70% 0px' });

sections.forEach(function (s) { observer.observe(s); });
