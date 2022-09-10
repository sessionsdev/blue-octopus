<script>
    export let bulmaSizeClass;
</script>
<section id="parallax" class="hero {bulmaSizeClass}">
    <div class="hero-body p-2">
      <div class="container">
          <slot></slot>
      </div>
    </div>
</section>

<style>
    #parallax  {
        background: var(--background);
        background-attachment: fixed;
        background-position: center;
        background-repeat: no-repeat;
        background-size: cover;
        border-top: 5px solid black;
}
</style>