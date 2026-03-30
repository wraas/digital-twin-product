/* W.R.A.A.S. — We Rickroll Absolutely Anyone, Seriously. */
console.log('%c W.R.A.A.S. ', 'background:#00c8ff;color:#050d1a;font-weight:900;font-size:1.4rem;padding:4px 14px;border-radius:3px;letter-spacing:0.1em;');
console.log('%c We Rickroll Absolutely Anyone, Seriously. ', 'color:#0065b3;font-size:1rem;font-weight:700;');
console.log('%c Happy April Fools 2026 ', 'color:#c47d00;font-size:0.9rem;font-style:italic;');
console.log('%c https://wraas.click ', 'color:#666;font-size:0.85rem;text-decoration:underline;');
console.log('%c You know the rules. ', 'color:#999;font-size:0.8rem;font-style:italic;');

// Mark that the user has read at least one docs page
try { localStorage.setItem('wraas-docs-read', '1'); } catch (e) {}

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

// Pagefind search
(function () {
  var searchBox = document.querySelector('.search-box');
  if (!searchBox) return;

  // Create results dropdown
  var results = document.createElement('div');
  results.className = 'search-results';
  results.setAttribute('role', 'listbox');
  results.setAttribute('aria-label', 'Search results');
  searchBox.parentNode.style.position = 'relative';
  searchBox.parentNode.appendChild(results);
  searchBox.setAttribute('role', 'combobox');
  searchBox.setAttribute('aria-expanded', 'false');
  searchBox.setAttribute('aria-controls', 'search-results');
  results.id = 'search-results';

  var pagefind = null;
  var debounce = null;

  async function initPagefind() {
    if (pagefind) return pagefind;
    try {
      // Resolve path relative to docs root
      var base = document.querySelector('link[rel="stylesheet"][href*="docs.css"]');
      var pagefindPath = base ? base.href.replace(/assets\/docs\.css$/, '_pagefind/pagefind.js') : '/_pagefind/pagefind.js';
      pagefind = await import(pagefindPath);
      await pagefind.init();
    } catch (e) {
      // Pagefind not available (local dev without index)
      pagefind = null;
    }
    return pagefind;
  }

  async function doSearch(query) {
    var pf = await initPagefind();
    if (!pf) {
      results.innerHTML = '<div class="search-result-item search-result-empty">Search index not available. Run: just index</div>';
      results.style.display = 'block';
      searchBox.setAttribute('aria-expanded', 'true');
      return;
    }
    var search = await pf.search(query);
    var items = await Promise.all(search.results.slice(0, 8).map(function (r) { return r.data(); }));

    if (items.length === 0) {
      results.innerHTML = '<div class="search-result-item search-result-empty">No results found.</div>';
    } else {
      results.innerHTML = items.map(function (item) {
        return '<a class="search-result-item" href="' + item.url + '" role="option">' +
          '<div class="search-result-title">' + item.meta.title + '</div>' +
          '<div class="search-result-excerpt">' + item.excerpt + '</div>' +
          '</a>';
      }).join('');
    }
    results.style.display = 'block';
    searchBox.setAttribute('aria-expanded', 'true');
    activeIndex = -1;
    searchBox.removeAttribute('aria-activedescendant');
  }

  searchBox.addEventListener('input', function () {
    clearTimeout(debounce);
    var query = searchBox.value.trim();
    if (query.length < 2) {
      results.style.display = 'none';
      searchBox.setAttribute('aria-expanded', 'false');
      return;
    }
    debounce = setTimeout(function () { doSearch(query); }, 200);
  });

  // Close on click outside
  document.addEventListener('click', function (e) {
    if (!searchBox.contains(e.target) && !results.contains(e.target)) {
      results.style.display = 'none';
      searchBox.setAttribute('aria-expanded', 'false');
    }
  });

  // Keyboard navigation
  var activeIndex = -1;

  function setActiveResult(index) {
    var items = results.querySelectorAll('.search-result-item[role="option"]');
    items.forEach(function (el) { el.classList.remove('is-active'); el.removeAttribute('id'); });
    activeIndex = index;
    if (index >= 0 && index < items.length) {
      items[index].classList.add('is-active');
      items[index].id = 'search-active-option';
      searchBox.setAttribute('aria-activedescendant', 'search-active-option');
      items[index].scrollIntoView({ block: 'nearest' });
    } else {
      searchBox.removeAttribute('aria-activedescendant');
    }
  }

  searchBox.addEventListener('keydown', function (e) {
    var items = results.querySelectorAll('.search-result-item[role="option"]');
    if (e.key === 'Escape') {
      results.style.display = 'none';
      searchBox.setAttribute('aria-expanded', 'false');
      activeIndex = -1;
      searchBox.removeAttribute('aria-activedescendant');
      searchBox.blur();
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (items.length > 0) setActiveResult(Math.min(activeIndex + 1, items.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (items.length > 0) setActiveResult(Math.max(activeIndex - 1, 0));
    } else if (e.key === 'Enter' && activeIndex >= 0 && activeIndex < items.length) {
      e.preventDefault();
      items[activeIndex].click();
    }
  });
})();

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
