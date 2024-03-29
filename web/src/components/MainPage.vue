<template>
  <v-container fluid mt-5 >
    <v-row>
      <!-- Use two empty columns on each side as dynamic gutters, I want gutters on a big
        screen but not on a small one -->
      <v-col cols=0 md=2></v-col>
      <!-- All content should go in here -->
      <v-col cols=12 md=8>
        <div class="mb-16 pb-10">
          <h1 class="primary--text">Soundcloud Shuffle</h1>
          <p>
            For those not aware, Soundcloud's desktop site does not properly shuffle
            all your liked songs; only the first few songs get shuffled. Enter the username
            from your user profile URL below to shuffle your entire library.
          </p>
        </div>

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

        <div v-else>
          <v-form
            ref="form"
            @submit.prevent="handleSubmit"
          >
            <div class="d-flex align-baseline flex-wrap">
              <!-- need the enter event so we blur event to trigger validation before we submit -->
              <v-text-field
                v-model="userInput"
                ref="formText"
                validate-on-blur
                :rules="rules"
                class="xs3 mr-5"
                style="flex-basis: 85%"
                label="Enter your Soundcloud username from your user profile URL (ex: test-user-25)"
                color="info"
                @keyup.native.enter="$refs.formText.blur()"
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

          <v-dialog
            v-model="showHelp"
            width="90%"
          >
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                color="secondary"
                small
                class="text-capitalize white--text mt-2"
                v-bind="attrs"
                v-on="on"
              >
                Need Help?
              </v-btn>
            </template>

            <v-card>
              <v-card-text>
                <v-img
                  src="/help.gif"
                  width="90%"
                />
              </v-card-text>

              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  color="primary"
                  text
                  @click="showHelp = false"
                >
                  Close
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>

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
        </div>

      </v-col>
      <v-col cols=0 md=2></v-col>
    </v-row>

    <v-snackbar
      v-model="hasError"
      :timeout="5000"
      color="#d9534f"
      top
    >
        {{ errorText || "Something went wrong" }}
      <template v-slot:action="{ attrs }">
        <v-btn
          color="white"
          text
          v-bind="attrs"
          @click="hasError = false"
        >
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script>
import SoundcloudPlayer from './SoundcloudPlayer.vue';

export default {
  name: 'MainPage',
  components: {
    SoundcloudPlayer,
  },
  data() {
    return {
      currentIndex: 0,
      currentSong: null,
      errorText: '',
      hasError: false,
      isLoading: false,
      likes: [],
      rules: [
        (v) => !!v || 'You must enter a username',
        (v) => (!!v && this.isValidUsername(v)) || 'Usernames can only contain letters, numbers, underscores, and hypens',
      ],
      showHelp: false,
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
        const res = await fetch(`api/likes?${new URLSearchParams({ url })}`);
        if (!res.ok) {
          const err = await res.json();

          this.hasError = true;
          this.errorText = err;
          return;
        }

        this.likes = await res.json();

        this.currentIndex = 0;
        this.currentSong = this.likes[this.currentIndex];
      } catch (error) {
        this.hasError = true;
        this.errorText = 'Something went wrong';
      } finally {
        this.isLoading = false;
      }
    },
    handleSubmit() {
      const isValid = this.$refs.form.validate();
      if (isValid) {
        const url = `https://soundcloud.com/${this.userInput.trim()}`;

        this.getLikes(url);
        this.$refs.form.reset();
      }
    },
    isValidUsername(username) {
      // matches all alphanumeric, underscore, and hypen. I think that should cover all usernames?
      const regex = /^[\w-]+$/;

      return regex.test(username);
    },
  },
  computed: {
    showNext() {
      return this.currentIndex !== this.likes.length - 1;
    },
    showPrevious() {
      // don't want user to go below 0 for the index
      return this.currentIndex !== 0;
    },
    songUrl() {
      return `https://w.soundcloud.com/player/?url=${this.currentSong}`;
    },
  },
};
</script>
