<template>
  <div>
    <v-sheet width="400" class="mx-auto">
      <p class="header-text">Welcome to FlySpray</p>
      <p class="sub-header-text">
        Flyspray is an open-source web-based project management and bug tracking
        system. It provides a platform for teams to collaborate, track tasks,
        and manage software development projects. With features such as task
        assignment, progress tracking, user permissions, Flyspray streamlines
        project communication and helps teams identify and resolve software
        issues efficiently. It is a tool commonly used in software development
        and other project-based industries to enhance teamwork and project
        organization.
      </p>
      <div class="buttons">
        <v-btn type="submit" block class="mt-2 btn" @click="goToRegisterPage()"
          >Register</v-btn
        >
        <v-btn type="submit" block class="mt-2 btn" @click="goToLoginPage()"
          >Sign In</v-btn
        >
      </div>
    </v-sheet>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  methods: {
    goToRegisterPage() {
      this.$router.push("/signup");
    },
    goToLoginPage() {
      this.$router.push("/login");
    },
  },
});
</script>

<style scoped>
.header-text {
  color: #6945c4;
  text-align: center;
  font-family: Poppins;
  font-size: 24px;
  font-style: normal;
  font-weight: 700;
  line-height: normal;
  text-transform: capitalize;
  margin-bottom: 2rem;
}
.sub-header-text {
  font-size: 1rem;
  margin-bottom: 1rem;
  font-weight: 400;
  color: #525252;
  font-family: "Poppins", sans-serif;
  margin-top: 1rem;
}
.buttons {
  width: 70%;
  text-align: center;
  align-items: center;
  margin: auto;
  margin-top: 2rem;
}
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
  margin-bottom: 0.8rem;
}
</style>
