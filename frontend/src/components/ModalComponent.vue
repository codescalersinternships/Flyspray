<template>
  <v-dialog
    class="dialog"
    v-model="internalDialog"
    @update:modelValue="(val:boolean) => $emit('close-dialog', val)"
    theme="dark"
  >
    <v-card>
      <v-card-title class="headline">
        <slot name="title" class="title"></slot>
      </v-card-title>
      <v-card-text>
        <slot name="modal-body" class="modal-body"></slot>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>

        <slot name="custom-button"></slot>
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
      internalDialog: false as boolean,
    };
  },
  watch: {
    dialog(newVal: boolean) {
      this.internalDialog = newVal;
    },
  },
};
</script>

<style scoped>
@media (max-width: 1100px) {
  .dialog {
    max-width: 100% !important;
    width: 100% !important;
  }
}
.v-dialog > .v-overlay__content > .v-card > .v-card-text,
.v-dialog > .v-overlay__content > form > .v-card > .v-card-text {
  padding: 0px 50px 25px 50px;
}
.v-card .v-card-title {
  padding-top: 25px;
  padding-left: 50px;
}
.v-card-actions {
  padding: 0px 50px 25px 0px;
}
.modal-body {
  padding: 100px 100px 100px 100px;
}
.title {
  font-weight: 500;
}
.dialog {
  max-width: 50%;
}
</style>
