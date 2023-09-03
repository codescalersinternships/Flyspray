<template>
  <v-dialog
    v-model="internalDialog"
    @update:modelValue="dialogUpdated"
    max-width="50%"
  >
    <v-card>
      <v-card-title class="headline">
        <slot name="title" class="title">Modal Title</slot>
      </v-card-title>
      <v-card-text>
        <slot></slot>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <slot name="custom-button"></slot>
        <v-btn class="btn" @click="closeModal">Close</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
export default {
  props: {
    dialog: Boolean,
  },
  data() {
    return {
      internalDialog: false,
    };
  },
  watch: {
    dialog(newVal) {
      this.internalDialog = newVal;
    },
  },
  methods: {
    closeModal() {
      this.internalDialog = false;
      this.$emit("close-dialog", false);
    },
    dialogUpdated(value: boolean) {
      if (!value) {
        this.closeModal();
      }
    },
  },
};
</script>

<style scoped>
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  /* font-family: Roboto; */

  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
}
.title {
  font-weight: 500;
}
</style>
