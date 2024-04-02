<template>
  <div>
    <el-scrollbar height="460px">
      <el-form ref="patchFormRef" :model="form" :rules="rules" label-width="auto" style="max-width: 600px;">
        <el-form-item label="选择机器：" prop="ip">
          <el-select v-model="form.ip" placeholder="please select IP" @change="getEnvsByIp">
            <el-option v-for="item in ips" :key="item.ip" :label="item.ip" :value="item.ip" />
          </el-select>
        </el-form-item>
        <el-form-item label="内核版本：" prop="version">
          <el-select v-model="form.version" placeholder="please select kernel version" @change="changeKernelVersion">
            <el-option v-for="item in buildEnvs" :key="item.version" :label="item.version" :value="item"
              :disabled="!(item.rpm.debugInfo && item.rpm.srcRpm)">
              <span>{{ item.version }}</span>&nbsp;
              <span
                :style="{ color: item.rpm.debugInfo && item.rpm.srcRpm ? 'var(--el-color-success)' : 'var(--el-color-danger)' }">
                {{ item.rpm.debugInfo && item.rpm.srcRpm ? '支持' : '不支持' }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="补丁类型：" prop="patchType">
          <el-radio-group v-model="form.patchType">
            <el-radio value="ACC" @click="showTypeInput = false">ACC</el-radio>
            <el-radio value="SGL" @click="showTypeInput = false">SGL</el-radio>
            <el-radio value="其他" @click="handlePatchType">其他</el-radio>
          </el-radio-group>
          <el-input v-model="form.patchType" v-if="showTypeInput"></el-input>
        </el-form-item>
        <el-form-item label="版本号：" prop="patchVersion">
          <el-input v-model="form.patchVersion" />
        </el-form-item>
        <el-form-item label="版本Release：" prop="patchRelease">
          <el-input v-model="form.patchRelease" />
        </el-form-item>
        <el-form-item label="CVE编号：" prop="patchDescription">
          <el-input v-model="form.patchDescription" />
        </el-form-item>
        <el-form-item label="patch文件：" prop="patchs">
          <el-upload class="upload" v-model:file-list="form.patchs" multiple :auto-upload="false">
            <template #trigger>
              <el-button type="success" plain>选择文件</el-button>
            </template>
            <template #tip>
              <div class="el-upload__tip" style="color: var(--el-color-warning);">
                *请选择.patch格式类型的文件(支持多选)
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form></el-scrollbar>
    <div style="display: flex; justify-content: flex-end;">
      <el-button @click="handleClose">取消</el-button>
      <el-button @click="handleConfirm(patchFormRef)">确定</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { addPatch } from '@/api/patch';
import { getAllHost, getBuildEnv } from "@/api/host";
import { RespCodeOK } from '@/api/request';
import type { BuildEnv, PatchForm } from '@/types/patch';
import { type Host } from "@/types/host";
import { ElMessage } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus'

const patchFormRef = ref<FormInstance>()
const showTypeInput = ref(false);
const form = reactive<PatchForm>({
  ip: '',
  buildKernelSrc: '',
  buildDebugInfo: '',
  version: '',
  patchType: 'ACC',
  patchs: [] as any,
  patchDescription: '',
  patchRelease: '',
  patchVersion: ''
})
const rules = reactive<FormRules<PatchForm>>({
  ip: [{ required: true, message: 'please select ip', trigger: 'blur' }],
  version: [{ required: true, message: 'please select kernel version' }],
  patchType: [{ required: true, message: 'please select kernel version' }],
  patchVersion: [{ required: true, message: 'please input patch version' }],
  patchRelease: [{ required: true, message: 'please input patch release' }],
  patchs: [{ required: true, message: 'please select patch file' }],
})

const emits = defineEmits(["close", "update"])
const ips = ref<Host[]>([])
onMounted(() => {
  getAllHost({ paged: false }).then(res => {
    if (res.data.code === RespCodeOK) {
      ips.value = res.data.data;
    }
  })
})

// 根据ip获取打包环境
const buildEnvs = ref<BuildEnv[]>([]);
const getEnvsByIp = (ip: string) => {
  getBuildEnv({ ip }).then(res => {
    if (res.data.code === RespCodeOK) {
      buildEnvs.value = res.data.data;
    }
  })
}

// 变更内核版本信息
const changeKernelVersion = (kernelRpms: any) => {
  if (!kernelRpms) return;
  form.version = kernelRpms.version;
  form.buildKernelSrc = kernelRpms.rpm.srcRpm;
  form.buildDebugInfo = kernelRpms.rpm.debugInfo;
}

// 处理补丁类型输入
const handlePatchType = () => {
  form.patchType = '';
  showTypeInput.value = true;
}

// 关闭弹窗
const handleClose = () => {
  emits('close')
}

// 确定上传
const fileList = ref([] as any)
const formData = new FormData();
const handleConfirm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid, _fields) => {
    if (valid) {
      formData.append('ip', form['ip'])
      formData.append('buildKernel', form['buildKernelSrc'])
      formData.append('buildDebugInfo', form['buildDebugInfo'])
      formData.append('version', form['version'])
      formData.append('patchType', form['patchType'])
      formData.append('patchDescription', form['patchDescription'])
      formData.append('patchVersion', form['patchVersion'])
      formData.append('patchRelease', form['patchRelease'])
      form.patchs.forEach((file: any) => {
        formData.append('upload', file.raw);
      });
      addPatch(formData).then(res => {
        if (res.data.code === RespCodeOK) {
          emits('close');
          emits('update');
          ElMessage.success(res.data.msg);
        } else {
          emits('close');
          ElMessage.error(res.data.msg);
        }
      })
    } else {
      ElMessage.error('error submit!');
    }
  })
}


</script>

<style scoped>
.upload {
  width: 100%;
}
</style>