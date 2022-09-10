<script>
  import "../app.scss";

  import TopNavBar from "$lib/component/TopNavBar.svelte";
  import Footer from "$lib/component/Footer.svelte";
</script>
<TopNavBar />

<slot />

<Footer />