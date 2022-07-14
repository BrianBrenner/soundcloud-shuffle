<template>
  <v-app>
    <v-main>
      <MainPage/>
    </v-main>
    <v-footer>
      <v-row class="d-flex justify-center pb-8">
        <a
          href="https://github.com/DoubleB123/soundcloud-shuffle"
          target="_blank"
          style="text-decoration: none;"
        >
          <v-icon large>
            mdi-github
          </v-icon>
        </a>
      </v-row>
    </v-footer>
  </v-app>
</template>

<script>
import MainPage from './components/MainPage.vue';

export default {
  name: 'App',
  components: { MainPage },
};
</script>

<style lang="scss">
.v-main__wrap {
  background-color: #f2f2f2;
}
</style>
