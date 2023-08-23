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
              :rules="[validateEmail]"
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

export default defineComponent({
  components: {
    InputForm,
  },
  data() {
    return {
      email: "" as string,
    };
  },
  computed: {
    errorEmail(): boolean {
      const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      return !pattern.test(this.email);
    },
  },
  methods: {
    validateEmail(value: string) {
      if (!value) return "Email is required";
      const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      return pattern.test(value) || "Invalid email format";
    },
    submitForm() {
      console.log("Email:", this.email);
    },
  },
});
</script>

<style scoped>
.header-text {
  font-size: 2rem;
  font-weight: 200;
  margin-bottom: 1rem;
  color: #8e73d3;
  font-family: "Fira Sans Extra Condensed", sans-serif;
}
.signin-text {
  color: #7d7d7d;
  margin-bottom: 1rem;
  margin-top: 1rem;
}
.btn {
  background-color: #8473f3;
  color: #ffffff;
}
.container {
  text-align: left;
  margin-bottom: 1rem;
}
</style>
