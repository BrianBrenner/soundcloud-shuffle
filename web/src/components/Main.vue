<template>
  <v-container fluid mt-5>
    <div
      class="d-flex justify-center"
    >
      <div style="width: 75%">
        <div v-if="!isLoaded" class="d-flex justify-center">
          <v-progress-circular
            indeterminate
            color="rgb(204, 68, 0)"
            size="40"
          />
        </div>
        <div v-else width="80%">
          <player 
            :url="songUrl"
            @next="handleNextSong"
            @previous="handlePreviousSong"
          />
        </div>
      </div>
    </div>
  </v-container>
</template>

<script>
import Player from './Player';

export default {
  name: 'Main',
  components: {
    Player,
  },
  data() {
    return {
      // clientId: '95f22ed54a5c297b1c41f72d713623ef',
      clientId: 'cyfcuTc4OKTJ09j8IoqWUZQkZ7QFN3p8',
      currentIndex: 0,
      currentSong: null,
      isLoaded: false,
      likes: [],
      shuffledLikes: [],
    };
  },
  async mounted() {
    // try {
    //   await this.getFavorites('/users/57886071/likes', {
    //     limit: 200,
    //     linked_partitioning: true,
    //   });


    // } catch (error) {
    //   console.log(error);
    // }


    // this.isLoaded = true;
  },
  methods: {
    handleNextSong() {
      this.currentIndex += 1;
      this.currentSong = this.shuffledLikes[this.currentIndex];
    },
    handlePreviousSong() {
      this.currentIndex -= 1;
      this.currentSong = this.shuffledLikes[this.currentIndex];
    },
    shuffleArray(array) {
      const copy = array.slice();

      for (let i = copy.length - 1; i > 0; i--) {
          const j = Math.floor(Math.random() * (i + 1));
          [copy[i], copy[j]] = [copy[j], copy[i]];
      }
      
      return copy;
    },
  },
  computed: {
    songUrl() {
      return `https://w.soundcloud.com/player/?url=https%3A//api.soundcloud.com${this.currentSong}`
    }
  }
}
</script>
