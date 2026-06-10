<script setup lang="ts">
import { ref } from 'vue'

const message = ref('PocketUnzip MVP')
</script>

<template>
  <div id="app">
    <h1>{{ message }}</h1>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background: #0a0e1a;
  color: #f1f5f9;
}

#app {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}
</style>
