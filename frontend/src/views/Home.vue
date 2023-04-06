<template>
  <main>
    <header>
      <img
        alt="PostgreScrutiniser logo"
        class="logo"
        src="@/assets/full-logo.svg"
        width="175"
        height="175"
      />
    </header>

    <HelloWorld />
    <TheWelcome />
  </main>
</template>

<script setup lang="ts">
import HelloWorld from "../components/HelloWorld.vue";
import TheWelcome from "../components/TheWelcome.vue";
</script>
