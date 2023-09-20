<script setup>
import DittoForm from './components/DittoForm.vue'</script>

<template>
  <el-row :gutter="10" justify="center">
    <el-col :xs="24" :sm="18" :md="16" :lg="14" :xl="12" >
      <div class="grid-content ep-bg-purple" style="height: 100vh;">
        <DittoForm/>
      </div>
    </el-col>
  </el-row>
</template>