<template>
  <div>
    <v-btn @click="$refs.dialog = true">Show Terms and Conditions</v-btn>

    <!-- Embed the BasicModal component and customize its content -->
    <ModalComponent ref="termsModal">
      <template v-slot:title> Terms and Conditions </template>
      <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit...</p>
      <!-- Add more content here as needed -->
    </ModalComponent>
  </div>
</template>
<script lang="ts">
import ModalComponent from "./ModalComponent.vue";
export default {
  components: {
    ModalComponent,
  },
};
</script>

<style scoped>
/* Add your custom styles here */
</style>
