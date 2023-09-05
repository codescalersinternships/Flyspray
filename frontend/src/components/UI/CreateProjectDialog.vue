<template>
  <div>
    <ModalComponent
      :dialog="internalDialog"
      @close-dialog="$emit('close', false)"
    >
      <template v-slot:title>
        <span class="text-h5">Add Project</span>
      </template>

      <template v-slot:modal-body>
        <v-container>
          <v-form v-model="formValid">
            <v-row>
              <v-col cols="12">
                <v-text-field
                  v-model="formValues.projectName"
                  :rules="(rules.emptyValues as any)"
                  label="Project Name"
                  required
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="formValues.projectDesc"
                  :rules="(rules.emptyValues as any)"
                  label="description"
                ></v-textarea>
              </v-col>
            </v-row>
          </v-form>
        </v-container>
      </template>
      <template v-slot:custom-button>
        <v-btn variant="text" color="btn" @click="$emit('close', false)">
          Close
        </v-btn>
        <v-btn
          color="primary"
          variant="text"
          @click="$emit('close', false)"
          :disabled="!formValid"
        >
          Save
        </v-btn>
      </template>
    </ModalComponent>
  </div>
</template>

<script lang="ts">
import ModalComponent from "@/components/ModalComponent.vue";
export default {
  components: {
    ModalComponent,
  },
  props: {
    dialog: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      formValues: {
        projectName: "",
        projectDesc: "",
      },
      rules: {
        emptyValues: [(val: string) => val.length || "Please write anything"],
      },
      formValid: false,
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
