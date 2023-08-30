<template>
  <div class="row">
    <logo-image class="child1"></logo-image>
    <div class="child2">
      <home-not-loggedin-component
        v-if="loggedIn"
      ></home-not-loggedin-component>
    </div>
  </div>
</template>

<script lang="ts">
import HomeNotLoggedinComponent from "../../components/home/HomeNotLoggedinComponent.vue";
import LogoImage from "../../components/LogoImage.vue";
import { defineComponent } from "vue";

export default defineComponent({
  name: "LoginPage",
  components: {
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
    };
  },
});
</script>
<style scoped>
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
