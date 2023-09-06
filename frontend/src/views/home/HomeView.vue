<template>
  <div>
    <div class="row" v-if="!loggedIn">
      <logo-image class="child1"></logo-image>
      <div class="child2">
        <home-not-loggedin-component
          v-if="loggedIn"
        ></home-not-loggedin-component>
      </div>
    </div>
    <v-app v-else>
      <navbar-sidebar-component
        @home-is-active="(value) => (isHomeActive = value)"
        @project-is-active="(value) => (isProjectActive = value)"
        @members-is-active="(value) => (isMembersActive = value)"
        @bugs-is-active="(value) => (isBugsActive = value)"
      ></navbar-sidebar-component>
      <v-main class="main-page">
        <v-container class="h-100">
          <v-row class="h-100">
            <main-page-component
              :isBugsActive="computedIsBugsActive"
              :isHomeActive="computedIsHomeActive"
              :isMembersActive="computedIsMembersActive"
              :isProjectActive="computedIsProjectActive"
            ></main-page-component>
            <latest-changes-component></latest-changes-component>
          </v-row>
          <footer-component></footer-component>
        </v-container>
      </v-main>
    </v-app>
  </div>
</template>

<script lang="ts">
import NavbarSidebarComponent from "../../components/NavbarSidebarComponent.vue";
import MainPageComponent from "../../components/home/MainPageComponent.vue";
import FooterComponent from "../../components/FooterComponent.vue";
import LatestChangesComponent from "../../components/LatestChangesComponent.vue";

import HomeNotLoggedinComponent from "../../components/home/HomeNotLoggedinComponent.vue";
import LogoImage from "../../components/LogoImage.vue";

import { defineComponent } from "vue";
export default defineComponent({
  name: "HomeLoggedInPage",

  components: {
    NavbarSidebarComponent,
    MainPageComponent,
    FooterComponent,
    LatestChangesComponent,
    HomeNotLoggedinComponent,
    LogoImage,
  },
  beforeMount() {
    const loggedIn = true;
    localStorage.setItem("loggedIn", JSON.stringify(loggedIn));
    console.log(localStorage.getItem("loggedIn"));
  },
  data() {
    return {
      loggedIn: localStorage.getItem("loggedIn"),
      isHomeActive: false as boolean,
      isProjectActive: false as boolean,
      isMembersActive: false as boolean,
      isBugsActive: false as boolean,
    };
  },
  computed: {
    computedIsBugsActive(): boolean {
      return this.isBugsActive;
    },
    computedIsHomeActive(): boolean {
      return this.isHomeActive;
    },
    computedIsMembersActive(): boolean {
      return this.isMembersActive;
    },
    computedIsProjectActive(): boolean {
      return this.isProjectActive;
    },
  },
});
</script>

<style scoped>
.main-page {
  background: #161616;
  color: #fff;
  font-family: Roboto;
  font-style: normal;
  line-height: normal;
}
.row {
  display: flex;
  height: auto;
  min-height: 100vh;
}
.child1 {
  flex-basis: 50%;
  width: 50%;
  text-align: center;
  background-color: #8457f7;
}
.child2 {
  flex-basis: 50%;
  width: 50%;
  text-align: center;
  margin-top: 5rem;
}

@media (max-width: 1100px) {
  .child1 {
    display: none;
  }
  .child2 {
    flex-basis: 100%;
    width: auto;
    padding: 0;
  }
}
</style>
