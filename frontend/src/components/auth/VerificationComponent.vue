<template>
  <div class="login-box">
    <div class="header-box">
      <p class="header-text" id="login-card-title">Verify Your Email</p>
      <div class="sub-box">
        <p class="sub-header-text" id="verification-card-title">
          Please enter the 4-digit code sent to your email.
        </p>
      </div>
    </div>
    <div class="form-box">
      <input-form>
        <v-sheet width="300" class="mx-auto">
          <v-form ref="form" @submit.prevent="submitForm">
            <div class="verification-code-container">
              <v-text-field
                v-for="(digit, index) in code"
                :key="index"
                v-model="code[index]"
                outlined
                dense
                class="verification-code-input"
                :ref="`digit${index}`"
                @input="formatCode()"
                @keypress="validateInput($event, index)"
              ></v-text-field>
            </div>
            <div class="center">
              <p class="error" v-if="error">Invalid code format</p>
            </div>
            <v-btn type="submit" block class="mt-2 btn">Submit</v-btn>
          </v-form>
          <div v-if="showResend">
            <hr class="form-separator" />

            <button
              class="link"
              :class="countDown != 0 ? 'disable' : ''"
              @click="resendCode()"
            >
              Resend Code
            </button>
            <p class="count-down" v-if="countDown != 0">{{ countDown }}</p>
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
      code: [null, null, null, null] as (number | null)[],
      error: false as boolean,
      showResend: false as boolean,
      countDown: 90 as number,
    };
  },
  methods: {
    formatCode() {
      for (let index = 0; index < 3; index++) {
        if (this.code[index] !== null) {
          const digitInputArray = this.$refs[
            "digit" + (index + 1)
          ] as HTMLInputElement[];
          digitInputArray[0].focus();
        }
      }
    },
    validateInput(event: KeyboardEvent, index: number) {
      if (this.code[index]) {
        event.preventDefault();
      }
      if (!/^\d$/.test(event.key)) {
        event.preventDefault();
      }
    },
    resendCode() {
      if (this.countDown === 0) {
        console.log("resend code");
        this.countDown = 90;

        this.startCountDown();
      } else {
        console.log("not allowed to resend code");
      }
    },
    startCountDown() {
      const countdownInterval = setInterval(() => {
        if (this.countDown > 0) {
          this.countDown--;
        } else {
          clearInterval(countdownInterval);
        }
      }, 1000);
    },
    submitForm() {
      this.error = false;
      if (this.code.every((code: number | null) => code != null)) {
        console.log("Verification code is filled:", this.code);
        this.showResend = true;
        if (this.countDown == 0 || this.countDown == 90) {
          this.countDown = 90;
          this.startCountDown();
        }
      } else {
        this.error = true;
        console.log("Verification code is not filled:", this.code);
      }
    },
  },
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css?family=Inter");
.header-box {
  margin-top: 5rem;
  margin-bottom: 2rem;
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
.sub-box {
  text-align: center;
  width: 22rem;
  margin: auto;
  margin-top: 1.5rem;
}
.sub-header-text {
  font-size: 1rem;
  margin-bottom: 1rem;
  font-weight: 400;
  color: #525252;
  font-family: "Poppins", sans-serif;
}
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
}
.link {
  color: var(--main-button, #8457f7);
  font-family: Poppins;
  font-size: 1rem;
  font-style: normal;
  font-weight: 100;
  line-height: normal;
  text-decoration: none;
  display: inline;
}
.count-down {
  display: inline;
  margin-left: 1rem;
  font-family: Poppins;
  color: #8457f7;
}
.disable {
  color: gray;
  cursor: not-allowed;
}
.form-separator {
  border: none;
  border-top: 1px solid #baa2f9;
  margin: 22px 0;
}
.verification-code-container {
  display: flex;
  width: 13rem;
  margin: auto;
}

.verification-code-input {
  flex: 1;
  height: 55px;
  font-size: 24px;
  text-align: center;
  border: 1px solid #ccc;
  border-radius: 4px 4px 0px 0px;
  margin-right: 8px;
  border: none;
  background-color: rgba(186, 162, 249, 0.4);
  outline-color: #5038ed;
}
.verification-code-input:focus-within {
  border: none;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.5);
}
.v-input__details {
  display: none !important;
}
.center {
  margin: auto;
  text-align: start;
  align-items: start;
  margin-bottom: 2rem;
  margin-top: 0.5rem;
  width: 13rem;
}
.error {
  transition-duration: 150ms;
  font-weight: 100;
  font-family: "Poppins", sans-serif;
  color: #c31031;
  font-size: 0.9rem;
}
</style>
