<template>
  <div id="swagger-ui"></div>
</template>

<script setup lang="ts">
import SwaggerUI from "swagger-ui";
import "swagger-ui/dist/swagger-ui.css";
import axios from "axios";
import type { openapiSpec } from "@/types/spec";

const urls = [
  `http://${import.meta.env.VITE_BACKEND_HOST}/api/docs/auth`,
  `http://${import.meta.env.VITE_BACKEND_HOST}/api/docs/file`,
  `http://${import.meta.env.VITE_BACKEND_HOST}/api/docs/resource-config`,
];

// IIFE for displaying swagger docs
(async () => {
  const mergedSpec = await fetchAndMergeSpecs(urls);
  SwaggerUI({
    dom_id: "#swagger-ui",
    spec: mergedSpec,
  });
})();

// function for merging multiple openapi specifiations
function mergeSpecs(specs: Array<openapiSpec>) {
  const mergedSpec = {
    servers: [{ url: `http://${import.meta.env.VITE_BACKEND_HOST}/api` }],
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