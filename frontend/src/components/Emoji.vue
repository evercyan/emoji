<template>
  <el-container class="page-emoji" v-loading="loading" :style="areaStyle">
    <el-aside class="side-area" style="width: 220px" v-if="list.length > 0">
      <el-menu :default-active="code" @select="onSelect">
        <el-menu-item v-for="(item, key) in list" :key="key" :index="item.code">
          <i class="el-icon-s-operation"></i>
          <span slot="title">{{ item.title }}</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main class="side-area" v-if="code">
      <el-card>
        <el-steps :active="active" finish-status="success" align-center>
          <el-step title="选择模板"></el-step>
          <el-step title="输入台词"></el-step>
          <el-step title="生成表情"></el-step>
        </el-steps>
      </el-card>
      <el-row>
        <el-col :span="12" class="left-area">
          <el-card class="box-card" v-if="info.gif">
            <div slot="header" class="clearfix">
              <span>示例</span>
            </div>
            <el-image :src="info.gif"></el-image>
          </el-card>
          <el-card class="box-card" v-if="info.gifData">
            <div slot="header" class="clearfix">
              <span>生成</span>
              <el-button type="success" size="mini" @click="downloadGif()" class="float-right">下载</el-button>
            </div>
            <el-image :src="info.gifData"></el-image>
          </el-card>
        </el-col>
        <el-col :span="12" class="right-area">
          <el-card class="box-card" v-if="code">
            <div slot="header" class="clearfix">
              <span>台词</span>
              <el-button
                type="danger"
                size="mini"
                @click="textList = info.sentences_list.slice()"
                class="float-right"
              >使用默认台词</el-button>
            </div>
            <el-input
              v-model="textList[index]"
              :placeholder="info.sentences_list[index]"
              v-for="(value, index) in textList"
              :key="index"
              clearable
            >
              <template slot="prepend">第 {{ index + 1 }} 句</template>
            </el-input>
            <el-button
              type="primary"
              @click="buildGif()"
              class="btn-build"
              :disabled="buildDisabled()"
            >生成</el-button>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: "page-emoji",
  data() {
    return {
      code: "",
      list: [],
      map: {},
      info: {},
      textList: [],
      active: 1,
      loading: false,
      areaStyle: {
        height: "",
      },
    };
  },
  watch: {
    textList(textList) {
      var _this = this;
      if (_this.active == 3) {
        return;
      }
      if (
        textList.length == _this.info.sentences_list.length &&
        textList.indexOf("") === -1
      ) {
        _this.active = 2;
      } else {
        _this.active = 1;
      }
    },
  },
  mounted() {
    var _this = this;
    // 重置 area 高度
    window.addEventListener("resize", _this.refreshAreaHeight);
    _this.refreshAreaHeight();
    // 加载表情模板列表
    _this.onLoading(true);
    _this.wails("GetTplList", "", function (result) {
      _this.onLoading(false);
      // 模板列表
      let list = JSON.parse(result);
      _this.list = list;

      // 模板 map
      let map = {};
      for (var i = 0; i < list.length; i++) {
        map[list[i].code] = list[i];
      }
      _this.map = map;

      _this.code = list[0].code;
      _this.onSelect(_this.code);
    });
  },
  methods: {
    refreshAreaHeight: function () {
      this.areaStyle.height = window.innerHeight + "px";
    },
    onSelect: function (code) {
      console.log("onSelect code", code);
      var _this = this;
      _this.code = code;
      _this.info = _this.map[code];
      _this.info.gif = require("@/assets/images/" + code + ".gif");
      // 进度重置
      _this.active = 1;
      // 写入空台词
      let textList = [];
      for (var i = 0; i < _this.info.sentences_list.length; i++) {
        textList.push("");
      }
      _this.textList = textList;
    },
    buildDisabled: function () {
      return this.textList.indexOf("") !== -1;
    },
    buildGif: function () {
      var _this = this;
      let param = JSON.stringify({
        code: _this.code,
        text_list: _this.textList,
      });
      _this.onLoading(true);
      _this.wails("BuildGif", param, function (result) {
        _this.onLoading(false);
        _this.info.gifData = result.data;
        _this.info.gifPath = result.path;
        _this.active = 3;
      });
    },
    downloadGif: function () {
      var _this = this;
      _this.wails("DownloadGif", _this.info.gifPath, function (result) {
        _this.$message.success(result);
      });
    },
    onLoading: function (loading) {
      this.loading = loading;
      if (loading) {
        setTimeout(() => {
          this.loading = false;
        }, 10000);
      }
    },
  },
};
</script>