<template>
  <v-app>
    <navbar-sidebar-component></navbar-sidebar-component>
    <v-main class="main-page">
      <v-container class="h-100">
        <v-row class="h-100">
          <main-page-component></main-page-component>
          <latest-changes-component></latest-changes-component>
        </v-row>
        <footer-component></footer-component>
      </v-container>
    </v-main>
  </v-app>
</template>
<style scoped>
.main-page {
  background: #161616;
  color: #fff;
  font-family: Roboto;
  font-style: normal;
  line-height: normal;
}
</style>
<script lang="ts">
import NavbarSidebarComponent from "@/components/NavbarSidebarComponent.vue";
import MainPageComponent from "@/components/MainPageComponent.vue";
import FooterComponent from "@/components/FooterComponent.vue";
import LatestChangesComponent from "@/components/LatestChangesComponent.vue";
import { defineComponent } from "vue";
export default defineComponent({
  name: "HomeLoggedInPage",
  components: {
    NavbarSidebarComponent,
    MainPageComponent,
    FooterComponent,
    LatestChangesComponent,
  },
});
</script>
