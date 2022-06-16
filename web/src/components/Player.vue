<template>
  <div>
    <iframe
      :src="urlWithOptions"
      width="100%" 
      allow="autoplay"
    >
    </iframe>
    <v-row justify="center">
      <button @click="$emit('next')">
        Next
      </button>
      <button @click="$emit('previous')">
        Previous
      </button>
    </v-row>
  </div>
</template>

<script>
import Widget from 'soundcloud-widget';

export default {
  name: 'Player',
  props: {
    url: {
      type: String,
      default: '',
      options: {
        auto_play: true,
      },
    },
  },
  async mounted() {
    const iframeElement = document.querySelector('iframe');
    const widget = new Widget(iframeElement);

    widget.on(Widget.events.READY, () => {
      widget.setVolume(50);
    });

    widget.on(Widget.events.FINISH, () => {
      console.log('done')
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
