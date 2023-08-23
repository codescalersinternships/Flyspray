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
      <v-sheet width="300" class="mx-auto">
        <v-form @submit.prevent="submitForm">
          <v-text-field
            prepend-inner-icon="mdi-email"
            label="Email"
            v-model="email"
          ></v-text-field>
          <v-text-field
            type="password"
            prepend-inner-icon="mdi-lock"
            label="Password"
            v-model="password"
          ></v-text-field>
          <div class="forgot-password-container">
            <router-link to="/forget" class="link"
              >Forgot Password?</router-link
            >
          </div>
          <v-btn type="submit" block class="mt-2 btn" :disabled="disable"
            >Sign In</v-btn
          >
        </v-form>
        <hr class="form-separator" />

        <p class="signin-text">
          Don't have account? please
          <router-link to="/signup" class="link">Sign Up</router-link>
        </p>
      </v-sheet>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  data() {
    return {
      email: "" as string,
      password: "" as string,
      errorEmail: true as boolean,
      errorPassword: false as boolean,
      disable: true as boolean,
    };
  },
  watch: {
    email() {
      const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      this.errorEmail = !pattern.test(this.email);
      this.disable = this.errorEmail || this.errorPassword;
    },
    password() {
      this.errorPassword = this.password.length < 8;
      this.disable = this.errorEmail || this.errorPassword;
    },
  },
  methods: {
    submitForm() {
      if (!this.disable) {
        console.log("login");
      }
    },
  },
};
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
