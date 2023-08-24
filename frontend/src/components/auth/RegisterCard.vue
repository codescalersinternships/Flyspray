<template>
  <div class="register-box">
    <div class="header-box">
      <p class="header-text" id="register-card-title">Create Your Account</p>
    </div>
    <div class="form-box">
      <input-form>
        <v-sheet width="300" class="mx-auto">
          <v-form ref="form" @submit.prevent="submitForm">
            <v-text-field
              prepend-inner-icon="mdi-account-outline"
              label="Username"
              required
              v-model="username"
              :rules="[validateUsernameRule]"
            ></v-text-field>
            <v-text-field
              prepend-inner-icon="mdi-email-outline"
              label="Email"
              v-model="email"
              required
              :rules="[validateEmailRule]"
            ></v-text-field>

            <v-text-field
              type="password"
              prepend-inner-icon="mdi-lock-outline"
              label="Password"
              v-model="password"
              required
              :rules="[validatePasswordRule]"
            ></v-text-field>
            <v-text-field
              type="password"
              prepend-inner-icon="mdi-lock-outline"
              label="Confirm Password"
              required
              :rules="[confirmPasswordRule]"
            ></v-text-field>
            <v-checkbox
              type="checkbox"
              label="I agree to all the Terms and Privacy policy."
              style="font-size: 15px"
              v-model="termsCheck"
            >
            </v-checkbox>
            <v-btn
              type="submit"
              block
              class="mt-2 btn"
              :disabled="
                errorPassword || errorEmail || errorUsername || errorTermsCheck
              "
              >Sign Up</v-btn
            >
          </v-form>
          <hr class="form-separator" />

          <p class="signin-text">
            Already have an account?
            <router-link to="/login" class="link">Sign in</router-link>
          </p>
        </v-sheet></input-form
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import InputForm from "./InputForm.vue";
import {
  validateEmail,
  ValidationResult,
  validatePasswordRegister,
  validateUsername,
} from "../../utils/validations";

export default defineComponent({
  components: {
    InputForm,
  },
  data() {
    return {
      email: "" as string,
      password: "" as string,
      username: "" as string,
      termsCheck: false as boolean,
      emailValidationResult: {} as ValidationResult,
      passwordValidationResult: {} as ValidationResult,
      usernameValidationResult: {} as ValidationResult,
    };
  },
  computed: {
    errorEmail(): boolean {
      return !validateEmail(this.email).isValid;
    },
    errorPassword(): boolean {
      return !validatePasswordRegister(this.password).isValid;
    },
    errorUsername(): boolean {
      return !validateUsername(this.username).isValid;
    },
    errorTermsCheck(): boolean {
      return !this.termsCheck;
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
    validatePasswordRule() {
      return (value: string) => {
        this.passwordValidationResult = validatePasswordRegister(value);
        return (
          this.passwordValidationResult.isValid ||
          this.passwordValidationResult.errorMessage
        );
      };
    },
    confirmPasswordRule() {
      return (value: string) => {
        return this.password == value || "Passwords do not match";
      };
    },
    // validateTermsRule() {
    //   return (
    //     this.termsCheck || "You have to agree to the terms and privacy policy"
    //   );
    // },
    validateUsernameRule() {
      return (value: string) => {
        this.usernameValidationResult = validateUsername(value);
        return (
          this.usernameValidationResult.isValid ||
          this.usernameValidationResult.errorMessage
        );
      };
    },
  },
  methods: {
    submitForm() {
      if (!this.errorEmail && !this.errorPassword && !this.errorUsername) {
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
.register-text {
  color: #7d7d7d;
  margin-bottom: 1rem;
}
.btn {
  background-color: #8473f3;
  color: #ffffff;
}
.forgot-password-container {
  text-align: right;
  margin-bottom: 1rem;
}

.link {
  display: inline-block;
  text-decoration: none;
  color: #8e73d3;
  font-size: 0.8rem;
}
.form-separator {
  border: none;
  border-top: 1px solid #ccc;
  margin: 20px 0;
}
</style>
