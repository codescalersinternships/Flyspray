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
              :class="{ 'error-field': errorEmail && isEmailClicked }"
              prepend-inner-icon="mdi-account"
              label="Email"
              v-model="email"
              required
              @click="isEmailClicked = true"
              class="input-label"
              :rules="[validateEmailRule]"
            ></v-text-field>

            <v-btn type="submit" block class="mt-2 btn" :disabled="errorEmail"
              >Reset password</v-btn
            >
          </v-form>
          <div class="container">
            <p class="signin-text">
              We'll send you a link to reset your password
            </p>
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
      isEmailClicked: false as boolean,
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
  margin-top: 2rem;
  color: #525252;
  text-align: center;
  font-family: Poppins;
  font-size: 1rem;
  font-style: normal;
  line-height: normal;
  margin-bottom: 3rem;
}
.input-label {
  color: var(--white-gray, #8f8f8f);
  font-family: Poppins;
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 400;
  line-height: normal;
  border-radius: 0.25rem;
  width: 100%;
  flex-shrink: 0;
}
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
}
.container {
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
.error-field {
  margin-bottom: 2rem;
  text-align: start;
  align-items: start;
}
.v-input__details {
  display: none !important;
}
</style>
