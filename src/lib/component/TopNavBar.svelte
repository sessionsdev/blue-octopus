
<script>
    let burgerMenu = false;
</script>

<nav class="navbar is-primary is-spaced" aria-label="main navigation">
    <div class="navbar-brand">
      <a class="navbar-item is-size-4" href="/">
        SESSIONS.DEV
    </a>
  
      <button href="#/" class="navbar-burger" class:is-active={burgerMenu} aria-label="menu" aria-expanded="false" data-target="navbarMenu" 
      on:click="{() => burgerMenu = !burgerMenu}">
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
        <span aria-hidden="true"></span>
      </button>
    </div>
  
    <div id="navbarMenu" class="navbar-menu" class:is-active={burgerMenu}>
      <div class="navbar-end">
        <a href="/" class="navbar-item">
            Home
        </a>
        <a href="/my-work" class="navbar-item">
            My Work
        </a>
        <a href="/my-clients" class="navbar-item">
            Clients
        </a>

        <div class="navbar-item">
          <div class="buttons">
            <a href="/hire-me" class="button is-info">
              <strong>HIRE ME</strong>
            </a>
          </div>
        </div>
      </div>
    </div>
  </nav>