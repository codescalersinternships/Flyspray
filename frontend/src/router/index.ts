import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import LoginView from "../views/auth/LoginView.vue";
import ForgetPasswordView from "../views/auth/ForgetPasswordView.vue";
import RegisterView from "../views/auth/RegisterView.vue";
import RegisterVerificationView from "../views/auth/RegisterVerificationView.vue";
import HomeNotLoggedinView from "../views/home/HomeNotLoggedinView.vue";

const routes: Array<RouteRecordRaw> = [
  {
    path: "/login",
    name: "login",
    component: LoginView,
  },
  {
    path: "/forget-password",
    name: "ForgetPassowrd",
    component: ForgetPasswordView,
  },
  {
    path: "/signup",
    name: "signup",
    component: RegisterView,
  },
  {
    path: "/register-verification",
    name: "register-verification",
    component: RegisterVerificationView,
  },
  {
    path: "/home",
    name: "home",
    component: HomeNotLoggedinView,
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
