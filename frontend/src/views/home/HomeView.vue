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
      <navbar-sidebar-component></navbar-sidebar-component>
      <v-main class="main-page">
        <v-container class="h-100">
          <v-row class="h-100">
            <main-page-component
              @new-project="showNewProjectDialog = true"
            ></main-page-component>
            <latest-changes-component></latest-changes-component>
          </v-row>
          <footer-component></footer-component>
        </v-container>
      </v-main>
    </v-app>
    <CreateProjectDialog :dialog="showNewProjectDialog" @close="close" />
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
import CreateProjectDialog from "@/components/UI/CreateProjectDialog.vue";
export default defineComponent({
  name: "HomeLoggedInPage",
  components: {
    NavbarSidebarComponent,
    MainPageComponent,
    FooterComponent,
    LatestChangesComponent,
    HomeNotLoggedinComponent,
    LogoImage,
    CreateProjectDialog,
  },
  beforeMount() {
    const loggedIn = true;
    localStorage.setItem("loggedIn", JSON.stringify(loggedIn));
    console.log(localStorage.getItem("loggedIn"));
  },
  data() {
    return {
      loggedIn: localStorage.getItem("loggedIn"),
      showNewProjectDialog: false,
    };
  },
  methods: {
    close() {
      this.showNewProjectDialog = false;
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
