<script setup lang="ts">
import WelcomeItem from "./WelcomeItem.vue";
import { IconBxBookContent, IconBxWrench } from "@iconify-prerendered/vue-bx";
</script>

<template>
  <WelcomeItem>
    <template #icon>
      <IconBxBookContent />
    </template>
    <template #heading>Documentation</template>

    PostgreSQL's official documentation for all of the runtime configuration
    options
    <a
      href="https://www.postgresql.org/docs/current/runtime-config-resource.html"
      target="_blank"
      rel="noopener"
      >can be found here</a
    >
    In addition, all of the API used for this project can be found by navigating
    to the "API Documentation" page found in the sidebar.
  </WelcomeItem>

  <WelcomeItem>
    <template #icon>
      <IconBxWrench />
    </template>
    <template #heading>Tooling</template>

    <span></span>
    This project is written with Vuejs, Typescript, Golang. Served and bundled
    with
    <a
      href="https://vitejs.dev/guide/features.html"
      target="_blank"
      rel="noopener"
      >Vite.</a
    >
    The API itself was generated using OpenAPI 3.0 tools. For backend:
    <a
      href="https://github.com/deepmap/oapi-codegen"
      target="_blank"
      rel="noopener"
      >oapi-codegen</a
    >
    For frontend:
    <a
      href="https://openapi-generator.tech/docs/installation"
      target="_blank"
      rel="noopener"
      >openapi-generator</a
    >
    <br />

    Instructions on how to generate API code as well as other development
    related details are available in <code>README.md</code>.
  </WelcomeItem>
</template>
