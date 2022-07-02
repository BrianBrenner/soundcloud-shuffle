<template>
  <v-container fluid mt-5 fill-height>
    <v-row
      v-if="isLoading"
      class="d-flex justify-center"
      key="loading"
    >
      <v-progress-circular
        indeterminate
        color="rgb(204, 68, 0)"
        size="40"
      />
    </v-row>

    <v-row
      v-else
      key="content"
    >
      <!-- Use two empty columns on each side as dynamic gutters, I want gutters on a big
        screen but not on a small one -->
      <v-col cols=0 md=2></v-col>
      <!-- All content should go in here -->
      <v-col cols=12 md=8>
        <v-form
          ref="urlForm"
          @submit="handleSubmit"
        >
          <div class="d-flex align-baseline flex-wrap">
            <!-- need the enter event so we blur event to trigger validation before we submit -->
            <v-text-field
              v-model="userInput"
              ref="urlText"
              validate-on-blur
              :rules="rules"
              class="xs3 mr-5"
              style="flex-basis: 85%"
              label="Enter your Soundcloud user profile URL (ex: https://soundcloud.com/test-user-25)"
              color="info"
              @keyup.native.enter="$refs.urlText.blur()"
            />
            <v-btn
              color="primary"
              class="text-capitalize white--text xs3"
              type="submit"
            >
              Shuffle Songs
            </v-btn>
          </div>
        </v-form>
        <div class="mt-15">
          <div v-if="likes.length">
            <soundcloud-player 
              :show-next="showNext"
              :show-previous="showPrevious"
              :url="songUrl"
              @next="handleNextSong"
              @previous="handlePreviousSong"
            />
          </div>
        </div>
      </v-col>
      <v-col cols=0 md=2></v-col>
    </v-row>

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
      isLoading: false,
      likes: [],
      rules: [
        v => !!v || 'You must enter a URL',
        v => !!v && this.isValidUrl(v) || 'The entered URL is not valid'
      ],
      userInput: '',
    };
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
    async getLikes(url) {
      try {
        this.isLoading = true;
        const res = await fetch('api/likes?' + new URLSearchParams({ url })) 
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
    handleSubmit() {
      const isValid = this.$refs.urlForm.validate();
      if (isValid) {
        // validatoion sucks so also do it by hand
        let formattedUrl;
        try {
          formattedUrl = new URL(this.userInput.trim()) 
        } catch (error) {
          // TODO: handle error
        }

        this.getLikes(formattedUrl.href)
        this.$refs.urlForm.reset();
      }
    },
    isValidUrl(url) {
      try {
        new URL(url) 
      } catch (error) {
        return false 
      }

      return true;
    }
  },
  computed: {
    showNext() {
      return this.currentIndex !== this.likes.length - 1
    },
    showPrevious() {
      // don't want user to go below 0 for the index
      return this.currentIndex !== 0
    },
    songUrl() {
      return `https://w.soundcloud.com/player/?url=${this.currentSong}`
    }
  }
}
</script>
