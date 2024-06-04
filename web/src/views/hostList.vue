<template>
  <div class="hostList">
    <my-table ref="tuneRef" :get-data="getAllHost" v-model:expandData="expandHost">
      <template #listName>构建环境列表</template>
      <template #button_bar>
        <el-input v-model="inputIp" style="max-width: 600px" placeholder="添加机器ip" clearable @change="handleAddAgent">
          <template #append>
            <el-button :icon="Plus" @click="handleAddAgent">新增</el-button>
          </template>
        </el-input>
      </template>
      <el-table-column type="expand">
        <template #default="{ row }">
          <el-table :height="row.ip ? 200 : 200" :data="patchRpms"
            style="width:80%; padding-left: 40px; border-radius: 4px;">
            <el-table-column prop="version" label="内核版本" align="center" width="220" />
            <el-table-column label="内核源码包" align="center" width="200">
              <template #default="{ row }">
                <el-icon v-if="row.rpm.srcRpm" color="#67C23A"><Select /></el-icon>
                <el-icon v-else color="#F56C6C">
                  <SemiSelect />
                </el-icon>
              </template>
            </el-table-column>
            <el-table-column label="debug调试安装包" align="center" width="200">
              <template #default="{ row }">
                <el-icon v-if="row.rpm.debugInfo" color="#67C23A"><Select /></el-icon>
                <el-icon v-else color="#F56C6C">
                  <SemiSelect />
                </el-icon>
              </template>
            </el-table-column>
            <el-table-column label="是否支持热补丁制作" align="center" width="200">
              <template #default="{ row }">
                <el-tag type="success" v-if="row.rpm.srcRpm && row.rpm.debugInfo">支持</el-tag>
                <el-tag type="danger" v-else>不支持</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </template>
      </el-table-column>
      <el-table-column prop="ip" label="IP" align="center" width="140" />
      <el-table-column prop="osVersion" label="系统版本" align="center" />
      <el-table-column prop="kernelArch" label="系统架构" align="center" width="100" />
      <el-table-column prop="kernelVersion" label="内核版本" align="center" />
      <el-table-column prop="lastheartbeat" label="最近心跳" align="center" width="180" />
      <el-table-column prop="connect" label="连接状态" align="center">
        <template #default="{ row }">
          <div class="flex-center">
            <span class="status-dot" :style="{ backgroundColor: row.connect ? '#67C23A' : '#F56C6C' }"></span>
            <span :style="{ color: row.connect ? '#67C23A' : '#F56C6C', paddingLeft: 4 + 'px' }">
              {{ row.connect ? '连接' : '离线' }}
            </span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="180">
        <template #default="{ row }">
          <!-- <el-tooltip class="box-item" effect="light" content="内核源码包/debug调试包" placement="top">
            <el-button size="small" plain @click="showUpload(row)">上传</el-button>
          </el-tooltip> -->
          <el-popconfirm title="确定删除这台机器吗？" width="200" @confirm="handleDelete(row)">
            <template #reference>
              <el-button type="danger" size="small" plain>删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </my-table>
    <el-dialog title="上传补丁包" width="30%" v-model="showDialog" destroy-on-close @close="files = []">
      <el-upload class="upload" v-model:file-list="files" multiple :auto-upload="false">
        <template #trigger>
          <el-button type="success" plain>选择文件</el-button>
        </template>
        <template #tip>
          <div class="el-upload__tip">
            请选择rpm格式文件
          </div>
        </template>
      </el-upload>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showDialog = false" plain>取消</el-button>
          <el-button type="primary" @click="handleUpload">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, watchEffect } from "vue";
import { getAllHost, delHost, getBuildEnv, addHost, uploadRpm } from "@/api/host";
import { RespCodeOK } from "@/api/request";
import { Plus } from '@element-plus/icons-vue'
import { Host } from "@/types/host";
import { ElMessage } from "element-plus";
const tuneRef = ref();
const patchRpms = ref([] as any);
const expandHost = ref({} as Host);
// 获取打包环境
watchEffect(() => {
  if (expandHost.value && expandHost.value.ip)
    patchRpms.value = [];
  getBuildEnv({ ip: expandHost.value.ip }).then((res: any) => {
    if (res.data.code === RespCodeOK) {
      patchRpms.value = res.data.data;
    }
  })
})

const inputIp = ref<string>();
// 添加agent
const handleAddAgent = () => {
  if (inputIp.value) {
    addHost({ ip: inputIp.value }).then(res => {
      if (res.data.code === RespCodeOK) {
        handleRefresh();
        ElMessage.success(res.data.msg)
      } else {
        ElMessage.error(res.data.msg)
      }
    })
  }
}

// 上传
const showDialog = ref(false);
const files = ref([] as any);
const upload_ip = ref('');
/* const showUpload = (row: Host) => {
  if (!row.ip) return;
  upload_ip.value = row.ip;
  showDialog.value = true;

} */
const formData = new FormData();
const handleUpload = () => {
  formData.append('ip', upload_ip.value);
  files.value.forEach((file: any) => {
    formData.append('upload', file.raw);
  });
  uploadRpm(formData).then(res => {
    if (res.data.code === RespCodeOK) {
      handleRefresh();
      ElMessage.success(res.data.msg);
    } else {
      ElMessage.error(res.data.msg);
    }
  })
}

// 删除
const handleDelete = (row: Host) => {
  if (!row.ip) return;
  delHost(row.ip).then((res: any) => {
    if (res.data.code === RespCodeOK) {
      handleRefresh();
      ElMessage.success(res.data.msg);
    } else {
      ElMessage.error(res.data.msg);
    }
  })
  // tuneRef.value.handleDelete();
};
// 刷新
const handleRefresh = () => {
  tuneRef.value.handleRefresh();
}
</script>

<style lang="less" scoped>
.hostList {
  width: 98%;
  margin: 0 auto;
  height: calc(100% - 44px - 20px);
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-dot {
  display: inline-block;
  width: 7px;
  height: 7px;
  vertical-align: middle;
  border-radius: 50%;
}
</style>
