<template>
  <v-app>
    <v-main>
      <MainPage/>
    </v-main>
  </v-app>
</template>

<script>
import MainPage from './components/MainPage';

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