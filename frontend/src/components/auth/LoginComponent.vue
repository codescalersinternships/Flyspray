<template>
  <div class="login-box">
    <div class="header-box">
      <p class="header-text" id="login-card-title">Welcome Back to FLYSPRAY</p>
      <p class="sub-header-text" id="login-card-title">
        Web-based bug tracker system
      </p>
    </div>
    <div class="form-box">
      <p class="signin-text">Sign in to start managing your projects</p>
      <input-form>
        <v-sheet width="300" class="mx-auto">
          <v-form ref="form" @submit.prevent="submitForm">
            <v-text-field
              :class="{ 'error-field': errorEmail && isEmailClicked }"
              prepend-inner-icon="mdi-email-outline"
              label="Email"
              v-model="email"
              required
              class="input-label"
              @click="isEmailClicked = true"
              :rules="[validateEmailRule]"
            ></v-text-field>
            <v-text-field
              type="password"
              class="input-label"
              :class="{ 'error-field': errorPassword && isPasswordClicked }"
              prepend-inner-icon="mdi-lock-outline"
              label="Password"
              v-model="password"
              required
              @click="isPasswordClicked = true"
              :rules="[validatePasswordRule]"
            ></v-text-field>

            <div class="forgot-password-container">
              <router-link to="/forget-password" class="link"
                >Forgot Password?</router-link
              >
            </div>
            <v-btn
              type="submit"
              block
              class="mt-2 btn"
              :disabled="errorPassword || errorEmail"
              >Sign In</v-btn
            >
          </v-form>
          <hr class="form-separator" />

          <p class="signin-text">
            Don't have an account? please
            <router-link to="/signup" class="link signin-link"
              >Sign Up</router-link
            >
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
  validatePassword,
} from "../../utils/validations";

export default defineComponent({
  components: {
    InputForm,
  },
  data() {
    return {
      email: "" as string,
      password: "" as string,
      emailValidationResult: {} as ValidationResult,
      passwordValidationResult: {} as ValidationResult,
      isEmailClicked: false as boolean,
      isPasswordClicked: false as boolean,
    };
  },
  computed: {
    errorEmail(): boolean {
      return !validateEmail(this.email).isValid;
    },
    errorPassword(): boolean {
      return !validatePassword(this.password).isValid;
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
        this.passwordValidationResult = validatePassword(value);
        return (
          this.passwordValidationResult.isValid ||
          this.passwordValidationResult.errorMessage
        );
      };
    },
  },
  methods: {
    submitForm() {
      if (!this.errorEmail && !this.errorPassword) {
        console.log("Email:", this.email);
      } else {
        console.log("Form is not valid");
      }
    },
  },
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css?family=Inter");
.header-box {
  margin: 1rem;
}
.header-text {
  color: #6945c4;
  text-align: center;
  font-family: Poppins;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
  text-transform: capitalize;
}
.sub-header-text {
  font-size: 1rem;
  margin-bottom: 1rem;
  font-weight: 400;
  color: #525252;
  font-family: "Poppins", sans-serif;
}
.signin-text {
  margin-bottom: 1rem;
  color: #525252;
  text-align: center;
  font-family: Poppins;
  font-size: 1rem;
  font-style: normal;
  line-height: normal;
  margin-bottom: 3rem;
}
.input-label {
  color: var(--white-gray, #494747);
  font-family: Poppins;
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  border-radius: 0.25rem;
  background: rgba(240, 237, 255, 0.2);
  width: 100%;
  flex-shrink: 0;
}
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
}
.forgot-password-container {
  text-align: right;
  margin-bottom: 1rem;
}

.link {
  color: var(--main-button, #8457f7);
  font-family: Poppins;
  font-size: 1rem;
  font-style: normal;
  font-weight: 100;
  line-height: normal;
  text-decoration: none;
}
.signin-link {
  color: #8457f7;
  font-weight: 400;
}
.form-separator {
  border: none;
  border-top: 1px solid #baa2f9;
  margin: 20px 0;
}
.error-field {
  margin-bottom: 2rem;
  text-align: start;
  align-items: start;
}
.v-input__details {
  display: none !important;
}
</style>
