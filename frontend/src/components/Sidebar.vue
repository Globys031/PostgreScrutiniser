<template>
  <div class="sidebar" :class="{ open: isOpen, }">
    <div class="logo-details">
      <img class="icon" alt="logo" src="@/assets/mini-logo.svg" width="50" height="50" />
      <span class="logo_name">PostgreScrutiniser</span>
      <span class="svg-box">
        <IconBxMenu id="btn" v-if="!isOpen" @click=toggleCloseButton() />
        <IconBxMenuAltRight id="btn" v-else @click=toggleCloseButton() />
      </span>
    </div>
    <ul class="nav-list">
      <li>
        <a href="#">
          <span class="svg-box">
            <IconBxGridAlt />
          </span>
          <span class="links_name">Runtime Configurations</span>
        </a>
        <span class="tooltip">Dashboard</span>
      </li>
      <li>
        <a href="#">
          <span class="svg-box">
            <IconBxUser />
          </span>
          <span class="links_name">Placholder</span>
        </a>
        <span class="tooltip">Placholder</span>
      </li>
      <li>
        <a href="#">
          <span class="svg-box">
            <IconBxUser />
          </span>
          <span class="links_name">Placholder</span>
        </a>
        <span class="tooltip">Placholder</span>
      </li>
      <li class="profile">
        <div class="profile-details">
          <div class>
            <div class="username">username</div>
            <div class="domain">exampledomain.com</div>
          </div>
        </div>
        <!-- <span class="svg-box"> -->
        <IconBxLogOut id="log_out" @click="logout()" />
        <!-- </span> -->
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { IconBxMenu, IconBxMenuAltRight, IconBxGridAlt, IconBxUser, IconBxLogOut } from '@iconify-prerendered/vue-bx';

const isOpen = ref<boolean>(false)

function toggleCloseButton() {
  console.log('the heck');
  isOpen.value = !isOpen.value
}

function logout() {
  console.log('Logout not implemented yet');
}
</script>

<style scoped>
/* Google Font Link */
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@200;300;400;500;600;700&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: "Poppins", sans-serif;
}

.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  height: 100%;
  width: 78px;
  background: var(--vt-primary-background);
  padding: 6px 14px;
  z-index: 99;
  transition: all 0.5s ease;
  color: var(--vt-c-white);
}

.sidebar.open {
  width: 310px;
}

.sidebar .logo-details {
  height: 60px;
  display: flex;
  align-items: center;
}

.sidebar .logo-details .icon {
  opacity: 0;
  transition: all 0.5s ease;
}

.sidebar .logo-details .svg-box {
  display: flex;
  justify-content: center;
  position: absolute;
  right: 0;
}

.sidebar .logo-details .logo_name {
  color: #fff;
  font-size: 20px;
  font-weight: 600;
  opacity: 0;
  transition: all 0.5s ease;
}


.sidebar.open .logo-details .icon,
.sidebar.open .logo-details .logo_name {
  opacity: 1;
}

.sidebar .logo-details #btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  font-size: 22px;
  cursor: pointer;
  transition: all 0.5s ease;
}

.sidebar.open .logo-details #btn {
  text-align: right;
}

/* .sidebar .logo-details #btn {
  text-align: center;
} */

.sidebar .nav-list {
  margin-top: 20px;
  height: 100%;
}

.sidebar li {
  position: relative;
  margin: 8px 0;
  list-style: none;
}


.sidebar li .tooltip {
  position: absolute;
  top: -20px;
  left: calc(100% + 15px);
  z-index: 3;
  background: #fff;
  box-shadow: 0 5px 10px rgba(0, 0, 0, 0.3);
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 15px;
  font-weight: 400;
  opacity: 0;
  white-space: nowrap;
  pointer-events: none;
  transition: 0s;
}

.sidebar li:hover .tooltip {
  opacity: 1;
  pointer-events: auto;
  transition: all 0.4s ease;
  top: 50%;
  transform: translateY(-50%);
}

.sidebar.open li .tooltip {
  display: none;
}

.sidebar li a {
  display: flex;
  height: 100%;
  width: 100%;
  border-radius: 12px;
  align-items: center;
  text-decoration: none;
  transition: all 0.4s ease;
  color: var(--vt-c-white);

  /* 
  Any of these would work for wrapping text
  but then when sidebar is not expanded the view would be broken
  */
  /* overflow-wrap: break-word; */
  /* word-wrap: break-word; */
  /* flex-wrap: wrap; */
}

.sidebar li a:hover {
  background: #FFF;
}

.sidebar li a .links_name {
  font-size: 15px;
  font-weight: 400;
  white-space: nowrap;
  opacity: 0;
  /* 
  for now leaving pointer-events uncommented
  but if navigation with links does not work,
  this is the culprit
  */
  pointer-events: none;
  transition: 0.4s;
}

.sidebar.open li a .links_name {
  opacity: 1;
  pointer-events: auto;
}

.sidebar li a:hover .links_name,
.sidebar li .tooltip,
.sidebar li a:hover svg {
  transition: all 0.5s ease;
  color: #11101D;
}

.sidebar .svg-box {
  color: #fff;
  height: 60px;
  min-width: 50px;
  font-size: 28px;
  text-align: center;
  line-height: 60px;
}

.sidebar li .svg-box {
  height: 50px;
  line-height: 50px;
  font-size: 18px;
  border-radius: 12px;
}


/* prie sito sugrizt. Ciuju tiesiog logout mygtukas bus ir tiek */
.sidebar li.profile {
  position: fixed;
  height: 60px;
  width: 78px;
  left: 0;
  bottom: -8px;
  padding: 10px 14px;
  background: #000000;
  transition: all 0.5s ease;
  overflow: hidden;
}

.sidebar.open li.profile {
  width: 310px;
}

.sidebar li .profile-details {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
}

.sidebar li.profile .username,
.sidebar li.profile .domain {
  font-size: 15px;
  font-weight: 400;
  color: #fff;
  white-space: nowrap;
}

.sidebar li.profile .domain {
  font-size: 12px;
}

.sidebar .profile #log_out {
  position: absolute;
  top: 50%;
  right: 0;
  transform: translateY(-50%);
  background: #1d1b31;
  /* width: 18px;
  height: 18px;
  line-height: 60px;
  border-radius: 0px; */
  transition: all 0.5s ease;
}

.sidebar.open .profile #log_out {
  /* width: 50px; */
  background: none;
}

@media (max-width: 420px) {
  .sidebar li .tooltip {
    display: none;
  }
}
</style>