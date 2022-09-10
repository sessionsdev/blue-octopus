<script>
    let cardData = {
        "id": "R_kgDOHvRuyA",
        "name": "colang",
        "pushedAt": "2022-07-29T19:59:51Z",
        "updatedAt": "2022-07-29T19:59:54Z",
        "url": "https://github.com/sessionsdev/colang",
        "description": "A cooklang text parser written in golang"
    }

</script>


<div class="card shadow-lg is-cursor-pointer">
    <div class="card-image cover-image is-overflow-hidden">
        <figure class="image">
        <img src="https://source.unsplash.com/user/himiwaybikes" alt="Placeholder image">
        </figure>
    </div>
    <div class="card-content">
        <div class="content has-text-grey-light">
            <h3>{cardData.name}</h3>
            {cardData.description}
        </div>
    </div>
    <footer class="card-footer has-background-white-bis">
        <a href={cardData.url} class="card-footer-item p-5 has-text-grey is-uppercase is-text-wide-1">
            GitHub
        </a>
    </footer>
</div>

<style>
.cover-image {
    max-height: 330px;
    border-top-left-radius: 3px;
    border-top-right-radius: 3px;
}

.cover-image img {
    transition: transform 1s, filter 2s ease-in-out;
    -webkit-transition: transform 1s, filter 2s ease-in-out;
    transform: scale(1.1);
    -webkit-transform: scale(1.1);
}

.card:hover .cover-image img {
    cursor: pointer;
    transform: scale(1.2);
    -webkit-transform: scale(1.2);
}

.is-overflow-hidden {
    overflow: hidden !important;
}

.is-cursor-pointer {
    cursor: pointer !important;
}

.shadow-lg {
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05) !important;
}
</style>