<template>
  <div>
    <iframe
      :src="urlWithOptions"
      width="100%"
      allow="autoplay"
      title="soundcloud-player"
    >
    </iframe>
    <v-row
      justify="center"
      class="mt-6"
    >
      <v-col cols="4"/>
      <v-col
        cols="4"
        class="d-flex justify-center"
      >
        <v-btn
          :disabled="!showPrevious"
          class="mr-2 text-capitalize white--text"
          color="primary"
          @click="$emit('previous')"
        >
          Previous
          <v-icon>
            mdi-skip-previous
          </v-icon>
        </v-btn>

        <v-btn
          :disabled="!showNext"
          class="ml-2 text-capitalize white--text"
          color="primary"
          @click="$emit('next')"
        >
          Next
          <v-icon>
            mdi-skip-next
          </v-icon>
        </v-btn>
      </v-col>
      <v-col cols="2"/>
      <v-col
        cols="2"
        align-self="center"
      >
        <v-slider
          v-model="volume"
          max="100"
          min="0"
          class="mt-1"
          prepend-icon="mdi-volume-high"
        />
      </v-col>

    </v-row>

  </div>
</template>

<script>
import _ from 'lodash';
import Widget from 'soundcloud-widget';

export default {
  name: 'SoundcloudPlayer',
  props: {
    showPrevious: {
      type: Boolean,
      required: true,
    },
    showNext: {
      type: Boolean,
      required: true,
    },
    url: {
      type: String,
      default: '',
    },
  },
  watch: {
    url() {
      if (this.widget) {
        // Just changing the src in the iframe isn't enough, we need to call this,
        // otherwise after the first song finishes or next is clicked, a new iframe
        // is loaded and all the events that are binded are unbinded. widget.load
        // keeps the events binded
        this.widget.load(this.url, { auto_play: true })
          .then(() => this.setVolume());
      }
    },
    volume() {
      this.volumeDebounced();
    },
  },
  data() {
    return {
      volume: 100,
      widget: null,
    };
  },
  async mounted() {
    const iframeElement = document.querySelector('iframe');
    this.widget = new Widget(iframeElement);

    this.widget.on(Widget.events.FINISH, () => {
      this.$emit('next');
    });
  },
  created() {
    this.volumeDebounced = _.debounce(this.setVolume, 10);
  },
  methods: {
    setVolume() {
      this.widget.setVolume(this.volume);
    },
  },
  computed: {
    urlWithOptions() {
      return `${this.url}&amp;auto_play=true`;
    },
  },
};
</script>
