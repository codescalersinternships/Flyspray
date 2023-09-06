<template>
  <div class="member-item-box">
    <div class="member-item-header">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        fill="currentColor"
        class="bi bi-columns-gap icon"
        viewBox="0 0 16 16"
      >
        <path
          d="M6 1v3H1V1h5zM1 0a1 1 0 0 0-1 1v3a1 1 0 0 0 1 1h5a1 1 0 0 0 1-1V1a1 1 0 0 0-1-1H1zm14 12v3h-5v-3h5zm-5-1a1 1 0 0 0-1 1v3a1 1 0 0 0 1 1h5a1 1 0 0 0 1-1v-3a1 1 0 0 0-1-1h-5zM6 8v7H1V8h5zM1 7a1 1 0 0 0-1 1v7a1 1 0 0 0 1 1h5a1 1 0 0 0 1-1V8a1 1 0 0 0-1-1H1zm14-6v7h-5V1h5zm-5-1a1 1 0 0 0-1 1v7a1 1 0 0 0 1 1h5a1 1 0 0 0 1-1V1a1 1 0 0 0-1-1h-5z"
        />
      </svg>
      <p id="project-name">project name</p>
    </div>
    <div class="member-item-details">
      <div class="member-in-project">
        <p>Member Name</p>
        <div v-if="isAdmin" class="icons-box">
          <div class="sub-box">
            <svg
              aria-hidden="true"
              height="16"
              viewBox="0 0 16 16"
              version="1.1"
              width="16"
              data-view-component="true"
              class="octicon octicon-shield-lock icon"
              fill="#FFF"
            >
              <path
                d="m8.533.133 5.25 1.68A1.75 1.75 0 0 1 15 3.48V7c0 1.566-.32 3.182-1.303 4.682-.983 1.498-2.585 2.813-5.032 3.855a1.697 1.697 0 0 1-1.33 0c-2.447-1.042-4.049-2.357-5.032-3.855C1.32 10.182 1 8.566 1 7V3.48a1.75 1.75 0 0 1 1.217-1.667l5.25-1.68a1.748 1.748 0 0 1 1.066 0Zm-.61 1.429.001.001-5.25 1.68a.251.251 0 0 0-.174.237V7c0 1.36.275 2.666 1.057 3.859.784 1.194 2.121 2.342 4.366 3.298a.196.196 0 0 0 .154 0c2.245-.957 3.582-2.103 4.366-3.297C13.225 9.666 13.5 8.358 13.5 7V3.48a.25.25 0 0 0-.174-.238l-5.25-1.68a.25.25 0 0 0-.153 0ZM9.5 6.5c0 .536-.286 1.032-.75 1.3v2.45a.75.75 0 0 1-1.5 0V7.8A1.5 1.5 0 1 1 9.5 6.5Z"
              ></path>
            </svg>
            <p>Admin</p>
          </div>
        </div>
      </div>
      <div class="member-in-project">
        <p>Member Name</p>
        <div v-if="isAdmin" class="icons-box">
          <div class="sub-box">
            <svg
              aria-hidden="true"
              height="16"
              viewBox="0 0 16 16"
              version="1.1"
              width="16"
              data-view-component="true"
              class="octicon octicon-shield-lock icon"
              fill="#FFF"
            >
              <path
                d="m8.533.133 5.25 1.68A1.75 1.75 0 0 1 15 3.48V7c0 1.566-.32 3.182-1.303 4.682-.983 1.498-2.585 2.813-5.032 3.855a1.697 1.697 0 0 1-1.33 0c-2.447-1.042-4.049-2.357-5.032-3.855C1.32 10.182 1 8.566 1 7V3.48a1.75 1.75 0 0 1 1.217-1.667l5.25-1.68a1.748 1.748 0 0 1 1.066 0Zm-.61 1.429.001.001-5.25 1.68a.251.251 0 0 0-.174.237V7c0 1.36.275 2.666 1.057 3.859.784 1.194 2.121 2.342 4.366 3.298a.196.196 0 0 0 .154 0c2.245-.957 3.582-2.103 4.366-3.297C13.225 9.666 13.5 8.358 13.5 7V3.48a.25.25 0 0 0-.174-.238l-5.25-1.68a.25.25 0 0 0-.153 0ZM9.5 6.5c0 .536-.286 1.032-.75 1.3v2.45a.75.75 0 0 1-1.5 0V7.8A1.5 1.5 0 1 1 9.5 6.5Z"
              ></path>
            </svg>
            <p>Admin</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  data() {
    return {
      isAdmin: true as boolean,
    };
  },
};
</script>

<style scoped>
.member-in-project {
  border-radius: 4px;
  background: #474157;
  width: 100%;
  display: flex;
  justify-content: space-between;
  padding: 12px 16px;
}
.member-item-box {
  border: 1px solid rgba(132, 87, 247, 1);
  margin-top: 2rem;
  border-radius: 4px;
}
.member-item-header {
  background-color: rgba(132, 87, 247, 1);
  color: white;
  font-size: 1.1rem;
  padding: 0.7rem;
  border-radius: 4px;
  display: flex;
  align-items: center;
}
.icon {
  margin-right: 0.5rem;
}
.member-item-details {
  padding: 0.5rem 2rem 1rem 2rem;
  font-size: 0.9rem;
  font-weight: 100;
  display: flex;
  flex-direction: column;
  gap: 12px;
  justify-content: flex-start;
  align-items: flex-start;
}
.icons-box {
  margin-top: 0.5rem;
  display: flex;
  justify-content: flex-start;
  gap: 1rem;
}
.icons-box .icon {
  margin-right: 0.1rem;
}
.bi-envelope {
  margin-right: 0.4rem !important;
}
.sub-box {
  display: flex;
}
@media (max-width: 600px) {
  .member-item-details {
    padding: 0;
    font-size: 0.7rem;
  }
}
</style>
