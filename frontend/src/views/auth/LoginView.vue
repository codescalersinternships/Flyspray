<template>
  <div class="row">
    <div class="child1">
      <v-img src="../../assets/logo.png" alt="logo" class="img"></v-img>
    </div>
    <div class="child2">
      <login-card></login-card>
    </div>
  </div>
</template>

<script lang="ts">
import LoginCard from "../../components/LoginCard.vue";
export default {
  name: "LoginPage",
  components: {
    LoginCard,
  },
};
</script>
<style scoped>
.row {
  display: flex;
  height: 100vh;
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
.img {
  width: 50%;
  margin: auto;
  margin-top: 30vh;
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
