<template>
  <v-row>
    <v-col cols="12" align-self="end">
      <v-footer class="footer" app>
        <div>
          <span>&copy; 2023 — FlySpray</span>
        </div>
      </v-footer>
    </v-col>
  </v-row>
</template>
<style scoped>
.footer {
  display: flex;
  padding: 16px 0px;
  justify-content: center;
  align-items: center;
  flex-shrink: 0;
  border-top: 1px solid #a9a9a9;
  background: #212121;
  color: #fff;
  text-align: center;
  font-family: Roboto;
  font-size: 14px;
  font-style: normal;
  font-weight: 500;
  line-height: normal;
}
</style>
<script lang="ts">
export default {
  data() {
    return {};
  },
};
</script>
