<template>
  <div id="content-container-center">
    <h3>a vs b!</h3>
    <form id="choice" name='form'>
      <button id="a" class="a" v-bind:style="{ opacity: current_choice == 'a' ? 0.5 : 1 }" value="a" v-on:click.stop.prevent="choice('a')">{{ option_a }}<i class="fa fa-check-circle" v-if="current_choice == 'a'"/></button>
      <button id="b" class="b" v-bind:style="{ opacity: current_choice == 'b' ? 0.5 : 1 }" value="b" v-on:click.stop.prevent="choice('b')">{{ option_b }}<i class="fa fa-check-circle" v-if="current_choice == 'b'"/></button>
    </form>
    <div id="tip">
      (Tip: you can change your vote)
    </div>
    <div v-if="current_choice">
      Current choice {{current_choice}}
    </div>
    <div id="hostname">
      Processed by container ID {{hostname}}
    </div>
  </div>
</template>

<script>
  export default {
    name: 'vote',
    data() {
      return {
        current_choice: null,
        option_a: 'Cat',
        option_b: 'Dog'
      }
    },
    methods: {
      choice: function (vote) {
        this.current_choice = vote
      }
    }
  }

</script>
