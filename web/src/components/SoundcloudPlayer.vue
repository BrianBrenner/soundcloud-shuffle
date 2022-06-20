<template>
  <div>
    <iframe
      :src="urlWithOptions"
      width="100%" 
      allow="autoplay"
    >
    </iframe>
    <v-row
      justify="center"
      class="mt-6"
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
    </v-row>
  </div>
</template>

<script>
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
  async mounted() {
    const iframeElement = document.querySelector('iframe');
    const widget = new Widget(iframeElement);

    widget.on(Widget.events.READY, () => {
      widget.setVolume(75);
    });

    widget.on(Widget.events.FINISH, () => {
      this.$emit('next');
    });
  },
  computed: {
    urlWithOptions() {
      return `${this.url}&amp;auto_play=true`;
    },
  },
}
</script>
