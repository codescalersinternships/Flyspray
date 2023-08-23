<template>
  <div class="login-box">
    <div class="header-box">
      <p class="header-text" id="login-card-title">
        It's okay! reset your password
      </p>
    </div>
    <div class="form-box">
      <input-form>
        <v-sheet width="300" class="mx-auto">
          <v-form ref="form" @submit.prevent="submitForm">
            <v-text-field
              prepend-inner-icon="mdi-account"
              label="Email"
              v-model="email"
              required
              :rules="[validateEmailRule]"
            ></v-text-field>

            <v-btn type="submit" block class="mt-2 btn" :disabled="errorEmail"
              >Continue</v-btn
            >
          </v-form>
          <div class="container">
            <p class="signin-text">We will send forget password link to you.</p>
          </div>
        </v-sheet></input-form
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import InputForm from "./InputForm.vue";
import { validateEmail, ValidationResult } from "../../utils/validations";

export default defineComponent({
  components: {
    InputForm,
  },
  data() {
    return {
      email: "" as string,
      emailValidationResult: {} as ValidationResult,
    };
  },
  computed: {
    errorEmail(): boolean {
      return !validateEmail(this.email).isValid;
    },
    validateEmailRule() {
      return (value: string) => {
        this.emailValidationResult = validateEmail(value);
        return (
          this.emailValidationResult.isValid ||
          this.emailValidationResult.errorMessage
        );
      };
    },
  },
  methods: {
    submitForm() {
      if (!this.errorEmail) {
        console.log("Email:", this.email);
      } else {
        console.log("Form is not valid");
      }
    },
  },
});
</script>

<style scoped>
.header-box {
  margin: 2rem;
}
.header-text {
  font-size: 2rem;
  font-weight: 200;
  margin-bottom: 1rem;
  color: #8e73d3;
  font-family: "Fira Sans Extra Condensed", sans-serif;
}
.sub-header-text {
  font-size: 1rem;
  margin-bottom: 1rem;
  font-weight: 400;
  color: #525252;
  font-family: "Poppins", sans-serif;
}
.signin-text {
  color: #7d7d7d;
  margin-bottom: 1rem;
}
.btn {
  background-color: #8473f3;
  color: #ffffff;
}
.container {
  text-align: right;
  margin-bottom: 1rem;
}
</style>
