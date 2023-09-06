<template>
  <v-col id="latest-changes-container" class="h-100" cols="3">
    <div id="latest-changes-innerdiv">
      <div id="latest-changes-title">Latest Changes</div>
      <div id="vertical-line">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="2"
          height="100%"
          viewBox="0 0 2 100%"
          fill="none"
        >
          <path d="M1 0L0.999982 404" stroke="#BAA2F9" />
        </svg>
        <div class="circle-container">
          <template v-for="index in 6" :key="index">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="15"
              height="15"
              viewBox="0 0 15 15"
              fill="none"
            >
              <circle cx="7.5" cy="7.5" r="7.5" fill="#BAA2F9" />
            </svg>
          </template>
        </div>
      </div>
    </div>
  </v-col>
</template>

<script lang="ts">
export default {
  data() {
    return {};
  },
};
</script>

<style scoped>
#latest-changes-container {
  display: flex;
  align-items: flex-start;
}
#latest-changes-innerdiv {
  width: 100%;
  height: 430px;
  border-radius: 4px;
  border: 1px solid var(--light-purple, #baa2f9);
  position: relative;
}
#latest-changes-title {
  text-align: start;
  padding-top: 16px;
  padding-left: 16px;
}
#vertical-line {
  position: absolute;
  left: 16px;
  padding-top: 30px;
  height: 80%;
}
.circle-container {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: absolute;
  top: 30px;
  left: -5px;
  height: calc(100% - 30px);
}

@media (max-width: 900px) {
  #latest-changes-container {
    display: none;
  }
}
</style>
