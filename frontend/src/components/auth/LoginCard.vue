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
              prepend-inner-icon="mdi-email"
              label="Email"
              v-model="email"
              required
              :rules="[validateEmail]"
            ></v-text-field>
            <v-text-field
              type="password"
              prepend-inner-icon="mdi-lock"
              label="Password"
              v-model="password"
              required
              :rules="[validatePassword]"
            ></v-text-field>

            <div class="forgot-password-container">
              <router-link to="/forget" class="link"
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
            Don't have account? please
            <router-link to="/signup" class="link">Sign Up</router-link>
          </p>
        </v-sheet></input-form
      >
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import InputForm from "./InputForm.vue";

export default defineComponent({
  components: {
    InputForm,
  },
  data() {
    return {
      email: "" as string,
      password: "" as string,
    };
  },
  computed: {
    errorEmail(): boolean {
      const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      return !pattern.test(this.email);
    },
    errorPassword(): boolean {
      return this.password.length < 8;
    },
  },
  methods: {
    validateEmail(value: string) {
      if (!value) return "Email is required";
      const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      return pattern.test(value) || "Invalid email format";
    },
    validatePassword(value: string) {
      if (!value) return "Password is required";
      return value.length >= 8 || "Password must be at least 8 characters long";
    },
    submitForm() {
      console.log("Email:", this.email);
      console.log("Password:", this.password);
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
