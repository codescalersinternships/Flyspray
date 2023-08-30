<template>
  <div class="container">
    <v-img src="../assets/404logo.png" id="err-logo"></v-img>
    <v-img src="../assets/404img.png" alt="logo" id="err-img"></v-img>
    <div class="row">
      <v-btn
        type="ReturnToHome"
        block
        class="mt-2 btn"
        @click="$router.push('/')"
        >Return To Home</v-btn
      >
      <v-btn
        type="ReturnToPrevious"
        block
        class="mt-2 btn"
        @click="$router.go(-1)"
        >Back</v-btn
      >
    </div>
  </div>
</template>

<style scoped>
.row {
  display: flex;
  margin-right: 220px;
  gap: 20px;
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
  margin-top: -195px !important;
  width: 70px;
}
.container {
  width: 50%;
  align-items: center;
  display: flex;
  flex-direction: column;
  margin-top: -166px;
}

#err-logo {
  height: auto;
  width: 20%;
  margin-bottom: -65%;
}

#err-img {
  height: auto;
  width: 80%;
}
@media (max-width: 1100px) {
  .container {
    width: 100%;
  }
  #err-logo {
    width: 30%;
  }
  #err-img {
    width: 100%;
    margin-top: -335px;
  }
  .row {
    margin-right: 170px;
  }
  .btn {
    width: 55px;
    margin-top: -250px !important;
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
