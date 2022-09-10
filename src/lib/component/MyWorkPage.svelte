<script>
    import GithubRepoCard from "./GithubRepoCard.svelte";
</script>

<h1 class="title">My Current Work</h1>
<div class="columns">
  <div class="section column is-one-third">
    <GithubRepoCard />
  </div>
</div>



<!-- 

query MyQuery {
    user(login: "sessionsdev") {
      bio
      repositories(
        first: 10
        privacy: PUBLIC
        orderBy: {field: PUSHED_AT, direction: DESC}
      ) {
        edges {
          node {
            id
            name
            pushedAt
            updatedAt
            url
            description
          }
        }
      }
    }
  }

   -->

<!-- 

   query MyQuery {
    user(login: "sessionsdev") {
      contributionsCollection {
        contributionCalendar {
          weeks {
            contributionDays {
              date
              contributionCount
              color
              weekday
            }
            firstDay
          }
          totalContributions
        }
      }
    }
  }

   -->