<template>
  <div>
    <v-app>
      <v-app-bar app class="navbar-container">
        <v-app-bar-nav-icon
          @click="drawer = !drawer"
          id="navbar-menu-icon"
        ></v-app-bar-nav-icon>

        <div class="navbar-left">
          <v-img id="navbar-logo" src="../assets/navbar-logo.png"></v-img>

          <v-img id="navbar-title" src="../assets/navbar-title.png"></v-img>
        </div>

        <div class="navbar-right">
          <v-img src="../assets/navbar-toggle-theme-icon.png"></v-img>
          <v-img src="../assets/navbar-user-icon.png"></v-img>
        </div>
      </v-app-bar>
    </v-app>
  </div>
</template>

<script lang="ts">
export default {
  data() {
    return {
      drawer: false,
    };
  },
};
</script>

<style scoped>
#navbar-menu-icon {
  color: #fbfbfb;
}
.navbar-container {
  background: var(--light-dark, #212121);
  box-shadow: 0px 3px 4px 0px rgba(0, 0, 0, 0.15);
}
.navbar-left,
.navbar-right {
  display: flex;
  align-items: center;
  height: 23px; /* Adjust the height as needed */
}
.navbar-right {
  margin-left: auto;
  width: 72px;
}
.navbar-left {
  width: 175px;
}
#navbar-title {
  width: 75%;
}
#navbar-logo {
  width: 20%;
}
</style>
