<template>
  <div class="container">
    <div id="err-logo">
      <v-img src="../assets/404logo.png"></v-img>
    </div>
    <div id="err-img">
      <v-img src="../assets/404img.png" alt="logo"></v-img>
    </div>
    <div class="row">
      <v-btn type="text" block class="mt-4 btn" @click="$router.push('/')"
        >Go Home</v-btn
      >
      <v-btn type="text" block class="mt-4 btn" @click="$router.go(-1)"
        >Back</v-btn
      >
    </div>
  </div>
</template>

<style scoped>
.row {
  width: 20%;
  display: flex;
  justify-content: center;
  gap: 10%;
}
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
  font-family: Inter;
  font-size: 0.75rem;
  font-style: normal;
  font-weight: 500;
  width: 100%;
}
.container {
  width: 100%;
  align-items: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

#err-logo {
  height: auto;
  width: 10%;
}

#err-img {
  height: auto;
  width: 50%;
}

@media (max-width: 400px) {
  .row {
    flex-direction: column;
    width: 100%;
  }
}
@media (max-width: 800px) {
  .row {
    flex-direction: column;
    width: 25%;
  }
}
</style>

<script lang="ts">
export default {
  data() {
    return {};
  },
};
</script>
