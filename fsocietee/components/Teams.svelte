<script>
let teams = [];
let maxTeamsFlag = false;

function addTeam() {
  if (teams.length >= 2) {
    return
  }
  teams = [...teams, { name: generateTeamName(), players: [''] }];
  if (teams.length >= 2) {
    maxTeamsFlag = true;
  };
};

function addMember(teamIndex) {
  teams[teamIndex].players = [...teams[teamIndex].players, ''];
  //TODO: check the max number of players in the room and do not allow 
  //to add if the total number is greater than the players
};

const nouns = [
  "Bananas",
  "Monkeys",
  "Nachos",
  "Giraffes",
  "Kittens",
  "Pickles",
  "Potatoes",
  "Wombats",
  "Noodles",
  "Unicorns",
  "Rainbows",
  "Cupcakes",
  "Bubbles",
  "Pandas",
  "Slippers",
  "Donuts",
  "Pillows",
  "Zombies",
  "Ninjas",
  "Pirates",
  "Avocados"
];

const adjectives = [
  "Confused",
  "Dizzy",
  "Sleepy",
  "Grumpy",
  "Sparkly",
  "Fluffy",
  "Slippery",
  "Bouncy",
  "Zippy",
  "Wiggly",
  "Giggly",
  "Chunky",
  "Funky",
  "Sneaky",
  "Cheeky",
  "Spicy",
  "Dorky",
  "Clumsy",
  "Crazy",
  "Goofy",
  "Hungry"
];

function generateTeamName() {
  const randomNoun = nouns[Math.floor(Math.random() * nouns.length)];
  const randomAdjective = adjectives[Math.floor(Math.random() * adjectives.length)];   

  return `${randomAdjective} ${randomNoun}`;
};

function validateTeams() {
  console.log(teams);
};
</script>


<div class="card mt-3">
  <div class="card-body">
    <h5 class="card-title">teams</h5>
    <div class="mt-3">
      {#each teams as team, i}
        <div>
          <h6> team {team.name} </h6>
          {#each team.players as player, j}
            <div class="form-group mt-1">
            <label for="team-{i}-player-{j}">player {j + 1}:</label>
            <input type="text" class="form-control" id="team-{i}-player-{j}" bind:value={teams[i].players[j]}>
            </div>
          {/each}
          <button type="button" class="btn btn-dark mt-1" on:click={() => addMember(i)}>add member</button>
          <div class="border-bottom border-2 border-dark my-4"></div>
        </div>        
      {/each}
      {#if !maxTeamsFlag}
      <button type="button" class="btn btn-dark mt-3" on:click={addTeam}>add team</button>
      {/if}
    </div>
    {#if maxTeamsFlag}
    <div class="mt-3">
      <button type="button" class="btn btn-dark mt-3" on:click={validateTeams}>validate teams</button>
    </div>
    {/if}
  </div>
</div>


