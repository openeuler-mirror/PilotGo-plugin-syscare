<template>
  <div class="patchList">
    <my-table ref="tuneRef" :get-data="getAllPatch">
      <template #listName>热补丁列表</template>
      <template #button_bar>
        <el-button type="primary" @click="showDialog = true">新增</el-button>
      </template>
      <el-table-column prop="ID" label="编号" width="60" align="center" />
      <el-table-column prop="ip" label="IP" width="130" align="center" />
      <el-table-column prop="buildVersion" label="构建版本" width="200" align="center" />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.status === '失败'" type="danger">{{ row.status }}</el-tag>
          <el-tag v-else-if="row.status === '完成'" type="success">{{ row.status }}</el-tag>
          <el-tag v-else type="primary">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="patchType" label="补丁类型" width="100" align="center" />
      <el-table-column prop="patchs" label="补丁文件" width="240" align="center">
        <template #default="{ row }">
          <p v-for="(item, index) in row.patchs.split(',')">{{ (index + 1) + '.' + item }}</p>
        </template>
      </el-table-column>
      <el-table-column prop="endTime" label="更新时间" align="center" />
      <el-table-column prop="hotPatch" label="热补丁包" align="center">
        <template #default="{ row }">
          {{ row.hotPatch }}
          <el-button type="primary" link :icon="Download" @click="downloadHotPatch(row)"></el-button>
        </template>
      </el-table-column>
      <el-table-column prop="patchKernel" label="热补丁内核源码包" align="center">
        <template #default="{ row }">
          {{ row.patchKernel }}
          <el-button type="primary" link :icon="Download" @click="downloadPatchKernel(row)"></el-button>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" align="center">
        <template #default="{ row }">
          <el-button size="small" plain @click="handleLog(row)">日志</el-button>
          <el-popconfirm title="确定删除这台机器吗？" width="200" @confirm="handleDelete(row)">
            <template #reference>
              <el-button type="danger" size="small" plain>删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </my-table>
    <el-dialog title="补丁包上传" width="40%" v-model="showDialog" destroy-on-close>
      <addPatch @close="showDialog = false" @update="handleRefresh" />
    </el-dialog>
    <el-drawer v-model="log_drawer" size="70%" title="构建日志" direction="rtl" destroy-on-close>
      <el-descriptions column="1">
        <el-descriptions-item label="日志状态：">{{ log_status }}</el-descriptions-item>
        <el-descriptions-item label="日志详情：">
          <pre>{{ log_detail }}</pre>
        </el-descriptions-item></el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { getAllPatch, delPatch } from "@/api/patch";
import { Patch } from "@/types/patch";
import addPatch from "./addPatch.vue";
import { RespCodeOK } from "@/api/request";
import { ElMessage } from "element-plus";
import { Download } from "@element-plus/icons-vue"
const tuneRef = ref();
const showDialog = ref(false);
const log_drawer = ref(false);

// 展示日志
const log_status = ref('');
const log_detail = ref('');
const handleLog = (row: Patch) => {
  log_drawer.value = true;
  log_status.value = row.stdout;
  log_detail.value = row.buildLog;
}

const closeDrawer = () => {
  log_detail.value = '';
}

// 下载热补丁包
const downloadHotPatch = (row: Patch) => {
  if (!row.hotPatch) return;
  window.open(window.location.origin + "/plugin/syscare/download/" + row.hotPatch + '?' + 'path=' + row.taskId);
}

// 下载热补丁内核源码包
const downloadPatchKernel = (row: Patch) => {
  if (!row.patchKernel) return;
  window.open(window.location.origin + "/plugin/syscare/download/" + row.patchKernel + '?' + 'path=' + row.taskId);
}

// 删除
const handleDelete = (row: Patch) => {
  if (!row.ID) return;
  delPatch({ id: row.ID }).then((res: any) => {
    if (res.data.code === RespCodeOK) {
      handleRefresh();
      ElMessage.success(res.data.msg);
    } else {
      ElMessage.error(res.data.msg);
    }
  })
};
// 刷新
const handleRefresh = () => {
  tuneRef.value.handleRefresh();
}
</script>

<style scoped>
.patchList {
  width: 98%;
  margin: 0 auto;
  height: calc(100% - 44px - 20px);
}
</style>
