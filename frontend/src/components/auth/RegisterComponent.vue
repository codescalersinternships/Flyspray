<template>
  <div class="register-box">
    <div class="header-box">
      <p class="header-text" id="register-component-title">
        Create Your Account
      </p>
    </div>
    <div>
      <input-form>
        <v-sheet width="300" class="mx-auto">
          <v-form class="form-box" ref="form" @submit.prevent="submitForm">
            <v-text-field
              class="input-label"
              :class="{ 'error-field': errorUsername && isUsernameClicked }"
              prepend-inner-icon="mdi-account-outline"
              label="Username"
              required
              v-model="username"
              @click="isUsernameClicked = true"
              :rules="[validateUsernameRule]"
            ></v-text-field>
            <v-text-field
              flat
              solo
              class="input-label"
              :class="{ 'error-field': errorEmail && isEmailClicked }"
              prepend-inner-icon="mdi-email-outline"
              label="Email"
              v-model="email"
              required
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
            <v-text-field
              type="password"
              class="input-label"
              :class="{ 'error-field': isConfirmPasswordClicked }"
              prepend-inner-icon="mdi-lock-outline"
              label="Confirm Password"
              required
              @click="isConfirmPasswordClicked = true"
              :rules="[confirmPasswordRule]"
            ></v-text-field>
            <v-checkbox type="checkbox" hide-details v-model="termsCheck">
              <template v-slot:label>
                <span class="terms-checkbox">
                  I agree to all the
                  <a
                    id="terms-link"
                    @click="
                      termsDialog = true;
                      $event.preventDefault();
                    "
                    >Terms</a
                  >
                  and
                  <a
                    id="privacy-policy-link"
                    @click="
                      privacyDialog = true;
                      $event.preventDefault();
                    "
                    >Privacy policy.</a
                  >
                </span>
                <ModalComponent
                  :dialog="termsDialog"
                  @close-dialog="(value) => (termsDialog = value)"
                >
                  <template v-slot:title>
                    <h3>Terms and Conditions</h3></template
                  >
                  <p class="bold">Last Updated: 09/03/2023</p>
                  <p class="bold">1. Acceptance of Terms</p>
                  <p>
                    By accessing or using Flyspray, you agree to comply with and
                    be bound by these terms and conditions. If you do not agree
                    with any part of these terms, please do not use Flyspray.
                  </p>

                  <p class="bold">2. License and Access</p>

                  <p class="bold">2.1. License:</p>
                  <p>
                    Subject to your compliance with these terms, Codescalers
                    grants you a limited, non-exclusive, non-transferable, and
                    revocable license to use Flyspray for your personal or
                    internal business purposes.
                  </p>

                  <template v-slot:custom-button>
                    <v-btn
                      class="btn"
                      @click="
                        termsCheck = true;
                        termsDialog = false;
                      "
                      >Accept Terms & Privacy Policy</v-btn
                    >
                  </template>
                </ModalComponent>
                <ModalComponent
                  :dialog="privacyDialog"
                  @close-dialog="(value) => (privacyDialog = value)"
                >
                  <template v-slot:title> <h3>Privacy Policy</h3> </template>
                  <p class="bold">Last Updated: 09/03/2023</p>
                  <p class="bold">1. Acceptance of Terms</p>
                  <p>
                    By accessing or using Flyspray, you agree to comply with and
                    be bound by these terms and conditions. If you do not agree
                    with any part of these terms, please do not use Flyspray.
                  </p>

                  <p class="bold">2. License and Access</p>

                  <p class="bold">2.1. License:</p>
                  <p>
                    Subject to your compliance with these terms, Codescalers
                    grants you a limited, non-exclusive, non-transferable, and
                    revocable license to use Flyspray for your personal or
                    internal business purposes.
                  </p>

                  <template v-slot:custom-button>
                    <v-btn
                      class="btn"
                      @click="
                        termsCheck = true;
                        privacyDialog = false;
                      "
                      >Accept Terms & Privacy Policy</v-btn
                    >
                  </template>
                </ModalComponent>
              </template>
            </v-checkbox>
            <v-btn
              type="Register"
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
            <router-link to="/login" id="signin-link">Sign in</router-link>
          </p>
        </v-sheet></input-form
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import InputForm from "./InputForm.vue";
import ModalComponent from "../ModalComponent.vue";
import {
  validateEmail,
  ValidationResult,
  validatePassword,
  validateUsername,
} from "../../utils/validations";

export default defineComponent({
  components: {
    InputForm,
    ModalComponent,
  },

  data() {
    return {
      email: "" as string,
      password: "" as string,
      username: "" as string,
      termsCheck: false as boolean,
      isUsernameClicked: false as boolean,
      isEmailClicked: false as boolean,
      isPasswordClicked: false as boolean,
      isConfirmPasswordClicked: false as boolean,
      emailValidationResult: {} as ValidationResult,
      passwordValidationResult: {} as ValidationResult,
      usernameValidationResult: {} as ValidationResult,
      termsDialog: false as boolean,
      privacyDialog: false as boolean,
    };
  },
  computed: {
    errorEmail(): boolean {
      return !validateEmail(this.email).isValid;
    },
    errorPassword(): boolean {
      return !validatePassword(this.password).isValid;
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
        this.passwordValidationResult = validatePassword(value);
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
@import url("https://fonts.googleapis.com/css?family=Inter");
.bold {
  font-weight: bold;
}
.header-box {
  margin: 2rem;
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
.register-text {
  color: #7d7d7d;
  margin-bottom: 1rem;
}
.signin-text {
  color: #525252;
  text-align: center;
  font-family: Poppins;
  font-size: 1rem;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
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
  height: 3rem;
  flex-shrink: 0;
}
.terms-checkbox {
  color: var(--body-text, #2d3748);
  font-family: Inter;
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 400;
  line-height: 140%; /* 1.05rem */
  letter-spacing: -0.015rem;
}
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
}
#terms-link,
#privacy-policy-link {
  color: #8457f7;
  font-family: Inter;
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 500;
  line-height: 140%;
  letter-spacing: -0.015rem;
}
.forgot-password-container {
  text-align: right;
  margin-bottom: 1rem;
}

#signin-link {
  color: var(--main-button, #9181f4);
  font-family: Poppins;
  font-size: 1rem;
  font-style: normal;
  font-weight: 600;
  line-height: normal;
  text-decoration: none;
}
.form-box {
  display: inline-flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 1rem;
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
