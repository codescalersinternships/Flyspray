<template>
  <v-col v-col id="main-section-container" class="h-100" cols="9" lg="9">
    <div id="main-section-innerdiv">
      <div class="header">
        <div class="sub-header">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="white"
          >
            <path
              d="M14 3V5H17.59L7.76 14.83L9.17 16.24L19 6.41V10H21V3M19 19H5V5H12V3H5C3.89 3 3 3.9 3 5V19C3 19.5304 3.21071 20.0391 3.58579 20.4142C3.96086 20.7893 4.46957 21 5 21H19C19.5304 21 20.0391 20.7893 20.4142 20.4142C20.7893 20.0391 21 19.5304 21 19V12H19V19Z"
            />
          </svg>
          <p class="header-text">Projects</p>
        </div>
        <div class="sub-header">
          <div class="search-box">
            <input placeholder="Search" class="search" variant="outlined" />
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              fill="currentColor"
              class="bi bi-search seacrh-icon"
              viewBox="0 0 16 16"
            >
              <path
                d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"
              />
            </svg>
          </div>
          <div>
            <v-btn
              class="btn"
              prepend-icon="mdi-plus"
              @click="$emit('new-project')"
            >
              Add Project
            </v-btn>
          </div>
        </div>
      </div>
      <div>
        <project-item></project-item>
        <project-item></project-item>
        <project-item></project-item>
        <project-item></project-item>
      </div>
    </div>
  </v-col>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import ProjectItem from "./ProjectItem.vue";

export default defineComponent({
  components: { ProjectItem },
  data() {
    return {};
  },
});
</script>
<style scoped>
.header {
  font-weight: 400;
  font-size: 1.3rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.header-text {
  margin-left: 0.5rem;
}
.sub-header {
  display: flex;
  align-items: center;
}
.search-box {
  display: flex;
  align-items: center;
}
.search {
  width: 12rem;
  outline-color: rgba(186, 162, 249, 1);
  color: rgba(186, 162, 249, 1);
  padding: 0.4rem;
  border: 1px solid rgba(186, 162, 249, 1);
  border-radius: 4px;
  font-size: 1rem;
}
.search:focus {
  outline: navajowhite;
}
::-webkit-input-placeholder {
  color: rgba(186, 162, 249, 1);
  font-size: 1rem;
}
.seacrh-icon {
  position: relative;
  right: 2rem;
  color: rgba(186, 162, 249, 1);
}
.btn {
  border-radius: 8px;
  background: linear-gradient(134deg, #9181f4 0%, #5038ed 100%);
  box-shadow: 0px 8px 21px 0px rgba(0, 0, 0, 0.16);
  color: #ffff;
}
#main-section-container {
  display: flex;
}

#main-section-innerdiv {
  width: 100%;
  border-radius: 4px;
  box-shadow: 0px 4px 4px 0px rgba(0, 0, 0, 0.25);
  /* padding: rem; */
}

@media (max-width: 900px) {
  #main-section-container {
    margin: auto;
  }
  .v-col-9 {
    flex: 0 0 100%;
    max-width: 100%;
  }
}
@media (max-width: 600px) {
  .search {
    width: 5rem;
    height: 2rem;
  }
  .btn {
    width: 6rem;
    font-size: 0.6rem;
  }
  .search::placeholder {
    color: transparent;
  }
  .v-col-9 {
    flex: 0 0 100%;
    max-width: 100%;
  }
}
</style>
