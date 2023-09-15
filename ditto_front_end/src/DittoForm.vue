<template>
  <el-form :model="form">
    <el-form-item>
      <div class="mt-4">
        <el-input
          v-model="form.url"
          placeholder="please input url"
          class="input-with-select"
          style="min-width: 800px"
        >
          <template #prepend>
            <el-select
              v-model="form.method"
              placeholder="Select"
              style="width: 115px"
            >
              <el-option
                v-for="item in methodOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </template>
          <template #append>
            <el-button @click="onSubmit" :icon="Promotion" />
          </template>
        </el-input>
      </div>
    </el-form-item>
    <el-row>
      <el-col :span="12">
        <el-form-item label="Loop">
          <el-input-number v-model="form.loop" controls-position="right" />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="Concurrency">
          <el-input-number
            v-model="form.concurrency"
            controls-position="right"
          />
        </el-form-item>
      </el-col>
    </el-row>
    <el-form-item>
      <el-tabs v-model="activeName" class="demo-tabs" @tab-click="handleClick">
        <el-tab-pane label="HEADER" name="first">
          <el-table :data="form.headers" style="width: 100%" max-height="250">
            <el-table-column label="key" width="150">
              <template #default="scope">
                <el-input v-model="scope.row.key"></el-input>
              </template>
            </el-table-column>
            <el-table-column label="value" width="150">
              <template #default="scope">
                <el-input v-model="scope.row.value"></el-input>
              </template>
            </el-table-column>
            <el-table-column label="process" width="120">
              <template #default="scope">
                <el-button
                  link
                  type="primary"
                  size="small"
                  @click.prevent="deleteRow(scope.row.id)"
                >
                  Remove
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button class="mt-4" style="width: 100%" @click="onAddItem">
            Add Item
          </el-button>
        </el-tab-pane>
        <el-tab-pane label="BODY" name="second" style="min-width: 500px">
          <el-input
            v-model="form.body"
            :autosize="{ minRows: 20, maxRows: 40 }"
            type="textarea"
            placeholder="Please input"
          />
        </el-tab-pane>
      </el-tabs>
    </el-form-item>
  </el-form>
</template>

<script lang="ts" setup>
import { reactive, ref } from "vue";
import axios from "axios";
import type { TabsPaneContext } from "element-plus";
import { Promotion } from "@element-plus/icons-vue";

const activeName = ref("first");

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event);
};

const id = ref(0);
const form = reactive({
  method: "POST",
  url: "",
  headers: [
    {
      id: id.value++,
      key: "",
      value: "",
    },
  ],
  body: "",
  loop: 1,
  concurrency: 1,
});

const onSubmit = () => {
  console.log(form);
  axios
    .post("http://localhost:8080/api/v1/ditto", {
      method: form.method,
      url: form.url,
      headers: form.headers,
      body: form.body,
      loop: form.loop,
      concurrency: form.concurrency,
    })
    .then((res) => {
      console.log(res);
    })
    .catch((err) => {
      console.log(err);
    });
};

const methodOptions = [
  {
    value: "1",
    label: "GET",
    disabled: true,
  },
  {
    value: "2",
    label: "POST",
    default: true,
  },
  {
    value: "3",
    label: "PUT",
    disabled: true,
  },
];

// headers
const deleteRow = (id: number) => {
  form.headers = form.headers.filter((item) => item.id !== id);
};

const onAddItem = () => {
  form.headers.push({
    id: id.value++,
    key: "",
    value: "",
  });
};
</script>