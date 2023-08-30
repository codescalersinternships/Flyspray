<template>
  <v-col id="main-section-container" class="h-100" cols="8">
    <div id="main-section-innerdiv">
      <h2>Main Section</h2>
    </div>
  </v-col>
</template>
<style scoped>
#main-section-container {
  display: flex;
}

#main-section-innerdiv {
  width: 100%;
  border-radius: 4px;
  border: 1px solid var(--main-color, #8457f7);
  box-shadow: 0px 4px 4px 0px rgba(0, 0, 0, 0.25);
}
</style>
<script lang="ts">
export default {
  data() {
    return {};
  },
};
</script>
