<template>
  <div class="my_table">
    <!-- table工具条 -->
    <el-row class="my_table_header">
      <div class="my_table_header_title">
        <slot name="listName"></slot>
        <div class="search">

        </div>
      </div>
      <div class="my_table_header_button">
        <el-input v-model="searchKey" style="max-width: 600px" placeholder="请输入关键字" clearable @change="getTableData">
          <template #append>
            <el-button :icon="Search" @click="getTableData" />
          </template>
        </el-input> &emsp;
        <slot name="button_bar"></slot>
      </div>
    </el-row>
    <!-- 列表 -->
    <div class="my_table_content" ref="tableBox">
      <el-table ref="myTableRef" :data="tableData" class="table" @select="handleRowSelectionChange"
        @selection-change="handleSelectinChange" :row-key="getRowKey" :expand-row-keys="expand_rowIds"
        @expand-change="onExpandChange">
        <slot></slot>
        <template #append>
          <slot name="append"></slot>
        </template>
        <template #empty>
          <el-empty description="无相关数据"></el-empty>
        </template>
      </el-table>
    </div>
    <!-- 分页 -->
    <div class="my_table_page">
      <el-pagination v-model:current-page="page.currentPage" v-model:page-size="page.pageSize"
        :page-sizes="[20, 25, 50, 75, 100]" :small="page.small" :background="page.background"
        layout="total, sizes, prev, pager, next, jumper" :total="page.total" @size-change="getTableData"
        @current-change="getTableData" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElTable, ElMessage, ElMessageBox } from "element-plus";
import { Search } from '@element-plus/icons-vue'
import { RespCodeOK, type ResultData } from "@/api/request";
const props = defineProps({
  getData: {
    type: Function,
    required: true,
  },
  delFunc: {
    type: Function,
    required: false,
  },
});
const emit = defineEmits(["handleSelect", "handleRowclick", "update:expandData"]);
const loading = ref(false);
const tableData = ref([] as any[]);
const myTableRef = ref<InstanceType<typeof ElTable>>();
const currentNum = ref(0); //复选框当前页数量
const selectedRows = ref([] as any[]);
const searchKey = ref(''); // 搜索框关键字
const page = reactive({
  total: 0,
  currentPage: 1,
  pageSize: 20,
  small: true,
  background: true,
});
onMounted(async () => {
  await getTableData();
});

// 获取表格数据
const getTableData = () => {
  loading.value = true;
  props.getData!({ page: page.currentPage, size: page.pageSize, paged: true, search: searchKey.value }).then(
    (res: { data: ResultData }) => {
      let result: ResultData = res.data;
      if (result && result.code === RespCodeOK) {
        loading.value = false;
        tableData.value = result.data!;
        currentNum.value = result.data!.length;
        page.total = Number(result.total);
      } else {
        loading.value = false;
        tableData.value = [];
        currentNum.value = 0;
        page.total = 10;
      }
    }
  );
};

// 绑定table的每一行的key
const getRowKey = (row: any) => {
  return row.id;
}
const expand_rowIds = ref<number[]>([]);
const expandId = ref<string>();
// 用户点击展开行只允许展开一行
const onExpandChange = (row: any, expandRows: any) => {
  expand_rowIds.value = [];
  if (row.id + '' !== expandId.value) {
    expandId.value = row.id + '';
    expand_rowIds.value.push(row.id);
  } else {
    expandId.value = '';
  }
  emit('update:expandData', row)
}

// 表格被选择数据发生变化
const handleSelectinChange = (rows: any[]) => {
  selectedRows.value = rows;
};

// 用户点击某一行的复选框
const handleRowSelectionChange = (val: [], _row: any[]) => {
  // 输出当前选中的所有行数组
  console.log(val);
};

// 刷新表格
const handleRefresh = () => {
  // 清空选项
  myTableRef.value?.clearSelection();
  // 重新获取数据
  getTableData();
};

// 删除
const handleDelete = () => {
  if (selectedRows.value.length == 0) return;
  ElMessageBox.confirm("确定要删除吗？", "提示", {
    type: "warning",
    confirmButtonText: "确定",
    cancelButtonText: "取消",
  }).then(() => {
    let ids = ref<number[]>([]);
    selectedRows.value.forEach((item) => {
      ids.value.push(item.id);
    });
    props.delFunc!({ ids: ids.value })
      .then((res: { data: ResultData }) => {
        if (res.data.code === RespCodeOK) {
          ElMessage.success(res.data.msg);
          handleRefresh();
        } else {
          ElMessage.error(res.data.msg);
        }
      })
      .catch((err: any) => {
        ElMessage.error("数据传输失败", err);
      });
  });
};

defineExpose({
  getTableData,
  handleDelete,
  handleRefresh,
});
</script>

<style lang="less" scoped>
.my_table {
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;
  width: 100%;
  height: 100%;

  &_header {
    width: 100%;
    height: 44px;
    padding: 0 6px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    &_title {
      font-size: 16px;
      width: 340px;
      display: flex;
      align-items: center;
      justify-content: space-between;

      .search {}
    }

    &_operation {
      flex: 2;
      display: flex;
      justify-content: flex-start;

      .operation {
        &-select {
          &-spanText {
            font-size: 12px;
            padding-left: 4px;
          }

          &-input {
            margin-left: 10px;
          }
        }
      }
    }

    &_button {
      display: flex;
      justify-content: flex-end;
    }
  }

  &_content {
    height: calc(100% - 64px - 30px - 20px);

    .table {
      width: 100%;
      height: 100%;
    }
  }

  &_page {
    height: 40px;
    padding: 0 4px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>@/types/host
