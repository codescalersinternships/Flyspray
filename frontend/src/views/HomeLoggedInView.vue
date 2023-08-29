<template>
  <div ud="navbar">
    <navbar-component></navbar-component>
  </div>
</template>
<script lang="ts">
import NavbarComponent from "@/components/NavbarComponent.vue";
import { defineComponent } from "vue";
export default defineComponent({
  name: "HomeLoggedInPage",
  components: {
    NavbarComponent,
  },
});
</script>
