<script>
    import TypewriterCodeBox from './TypewriterCodeBox.svelte'
    import ParallaxSection from './ParallaxSection.svelte';
</script>

    <!-- paralax section 1 -->
    <ParallaxSection --background="url('/img/main-background.png')" bulmaSizeClass="is-fullheight">
      <div class="columns mb-6 pb-6">
        <div class="column is-5 is-offset-7 is-hidden-mobile">
            <TypewriterCodeBox />
        </div>
      </div>
    </ParallaxSection>

    <!-- paralax section 2 -->
    <ParallaxSection --background="#4e95ec" bulmaSizeClass="is-halfheight">
      <div class="columns">
        <div class="column is-6">
          <h1 class="title is-1 ">Some Copy Here</h1>
          <!-- <hr class="content-divider"> -->
          <h2 class="subtitle">Lorem ipsum dolor, sit amet consectetur adipisicing elit. Doloremque veritatis sequi natus minima distinctio ullam deleniti quasi quisquam autem deserunt.</h2>
          <a href="/" class="button is-white is-inverted">Some Button</a>
        </div>
      </div>
    </ParallaxSection>

    <!-- paralax section 3 -->
    <ParallaxSection --background="url('/img/main-background.png')" bulmaSizeClass="is-fullheight">
      <div class="columns">
        <div class="column is-6 is-offset-6 box has-background-black">
          <h1 class="title is-1 ">SOME MORE COPY</h1>
          <hr class="content-divider">
          <h2 class="subtitle">Lorem ipsum dolor sit amet consectetur, adipisicing elit. Iure minus nam necessitatibus neque in perferendis eveniet dolorum assumenda dolores accusamus.</h2>
          <a href="/" class="button is-white">Another Button</a>
        </div>
      </div>
    </ParallaxSection>

<style>


.title {
  color: white;
  font-weight: 800;
  font-size: 4rem;
  text-shadow: 1px 1px 1px #000, 3px 3px 5px rgb(0, 0, 3);
}

.subtitle {
  color: white;
  font-size: 1.2rem;
  text-shadow: 1px 1px 1px #000, 3px 3px 5px rgb(0, 0, 3);
}

</style>