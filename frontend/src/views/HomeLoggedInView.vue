<template>
  <div ud="navbar">
    <navbar-sidebar-component></navbar-sidebar-component>
  </div>
</template>
<script lang="ts">
import NavbarSidebarComponent from "@/components/NavbarSidebarComponent.vue";
import { defineComponent } from "vue";
export default defineComponent({
  name: "HomeLoggedInPage",
  components: {
    NavbarSidebarComponent,
  },
});
</script>
