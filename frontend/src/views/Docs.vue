<template>
  <div id="swagger-ui"></div>
  <UiSpinner v-if="isLoading" />
</template>

<script setup lang="ts">
import { ref } from "vue";
import SwaggerUI from "swagger-ui";
import "swagger-ui/dist/swagger-ui.css";
import axios from "axios";
import { useSessionStore } from "@/stores/session";
import type { openapiSpec } from "@/types/spec";

const sessionStore = useSessionStore();

const urls = [
  `${sessionStore.baseAPIPath}/docs/auth`,
  `${sessionStore.baseAPIPath}/docs/file`,
  `${sessionStore.baseAPIPath}/docs/resource-config`,
];

const isLoading = ref<boolean>(true);

// IIFE for displaying swagger docs
(async () => {
  const mergedSpec = await fetchAndMergeSpecs(urls);
  SwaggerUI({
    dom_id: "#swagger-ui",
    spec: mergedSpec,
    requestInterceptor: (request: any) => {
      request.headers.Authorization = `Bearer ${sessionStore.token}`;
      return request;
    },
  });

  isLoading.value = false;
})();

// function for merging multiple openapi specifiations
function mergeSpecs(specs: Array<openapiSpec>) {
  const mergedSpec = {
    servers: [{ url: `${sessionStore.baseAPIPath}` }],
    components: {
      schemas: {},
    },
    info: {
      title: "API Documentation",
      version: "1.0.0",
    },
    openapi: "3.0.0",
    paths: {},
    security: {},
  };
  for (const spec of specs) {
    Object.assign(mergedSpec.paths, spec.paths);
    Object.assign(mergedSpec.components.schemas, spec.components.schemas);
    Object.assign(mergedSpec.security, spec.security);
  }

  return mergedSpec;
}

// Function for fetching openapi specifiation
async function fetchSpec(url: string) {
  try {
    const response = await axios.get(url);
    return response.data;
  } catch (error) {
    console.error(`Failed to fetch spec from ${url}`, error);
    return null;
  }
}

// Function for displaying different specifications as a singular spec
async function fetchAndMergeSpecs(urls: Array<string>) {
  const specsPromises = urls.map((url: string) => fetchSpec(url));
  const specs = await Promise.all(specsPromises);
  return mergeSpecs(specs);
}
</script>
