<template>
  <div class="row">
    <page-not-found-component></page-not-found-component>
  </div>
</template>
<style scoped>
.row {
  display: flex;
  height: 100vh;
  background: #161616;
  color: white;
  justify-content: center;
}
</style>
<script lang="ts">
import PageNotFoundComponent from "@/components/PageNotFoundComponent.vue";
import { defineComponent } from "vue";

export default defineComponent({
  name: "RegisterPage",
  components: {
    PageNotFoundComponent,
  },
});
</script>
