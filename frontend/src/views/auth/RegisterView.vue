<template>
  <div class="row">
    <logo-image class="child1"></logo-image>
    <div class="child2">
      <register-component></register-component>
    </div>
  </div>
</template>

<script lang="ts">
import LogoImage from "../../components/LogoImage.vue";
import RegisterComponent from "../../components/auth/RegisterComponent.vue";
import { defineComponent } from "vue";

export default defineComponent({
  name: "RegisterPage",
  components: {
    LogoImage,
    RegisterComponent,
  },
});
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
