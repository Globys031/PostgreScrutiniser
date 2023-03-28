<template>
  <main>
    <header>
      <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="125" height="125" />

      <div class="wrapper">
        <nav>
          <RouterLink to="/">Home</RouterLink>
          <RouterLink to="/about">About</RouterLink>
          <RouterLink to="/docs">Docs</RouterLink>
        </nav>
      </div>
    </header>

    <HelloWorld msg="You did it!" />
    <TheWelcome />
  </main>
</template>

<script setup lang="ts">
import HelloWorld from "../components/HelloWorld.vue";
import TheWelcome from "../components/TheWelcome.vue";
</script>