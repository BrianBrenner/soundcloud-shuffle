<template>
  <v-container fluid mt-5>
    <div
      class="d-flex justify-center"
    >
      <div style="width: 75%">
        <div v-if="isLoading" class="d-flex justify-center">
          <v-progress-circular
            indeterminate
            color="rgb(204, 68, 0)"
            size="40"
          />
        </div>
        <div v-else width="80%">
          <soundcloud-player 
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
import SoundcloudPlayer from './SoundcloudPlayer';

export default {
  name: 'MainPage',
  components: {
    SoundcloudPlayer,
  },
  data() {
    return {
      currentIndex: 0,
      currentSong: null,
      isLoading: true,
      likes: [],
    };
  },
  mounted() {
    this.getLikes();
  },
  methods: {
    handleNextSong() {
      this.currentIndex += 1;
      this.currentSong = this.likes[this.currentIndex];
    },
    handlePreviousSong() {
      this.currentIndex -= 1;
      this.currentSong = this.likes[this.currentIndex];
    },
    async getLikes() {
      try {
        this.isLoading = true;
        const res = await fetch('api/likes?' + new URLSearchParams({ url: 'https://soundcloud.com/brian-brenner-4' })) 
        this.likes = await res.json()

        this.currentIndex = 0
        this.currentSong = this.likes[this.currentIndex]
      } catch (error) {
        // TODO: handle errors
        console.log(error)
      } finally {
        this.isLoading = false;
      }
    },
  },
  computed: {
    songUrl() {
      return `https://w.soundcloud.com/player/?url=${this.currentSong}`
    }
  }
}
</script>
