# emoji

> 基于 go wails 和 vue2 实现的表情锅生成工具 mac app

---

#### QA

```
Q: 思路和模板资源文件来源
A: 看这里: https://github.com/xtyxtyx/sorry

Q: 核心依赖 ffmpeg
A: gif 生成依赖 ffmpeg, 需提前安装 /usr/local/bin/ffmpeg

Q: 遇到 "读取配置文件失败" 怎么解决? 
A: wails 打包无法包含资源文件, 故配置文件, 视频资源等均存于 github, 遇到此问题, 手动打开 https://raw.githubusercontent.com/evercyan/cantor/master/assets/emoji/template.json, 无法访问请参照: http://idayer.com/speed-github-githubusercontent-page-with-hosts/

Q: 系统运行日志
A: /tmp/emoji.log

Q: 支持哪些系统 
A: 仅 Mac 10.14+ 亲测
```

---

#### RUN

```sh
# 安装 ffmpeg
brew install ffmpeg
/usr/local/bin/ffmepg -version

# 安装 wails 
go get -u github.com/wailsapp/wails/cmd/wails
wails -help

# 下载 emoji
git clone https://github.com/evercyan/emoji

# 安装前端组件
cd ./emoji/frontend/
npm install

# 启动后端服务
cd ./emoji/
sh run.sh debug

# 启动前端服务
cd ./emoji/frontend
npm run serve

# 打开 http://127.0.0.1:8080/
```

```sh
# 生成可执行文件 ./build/emoji
sh run.sh test 

# 生成 mac app ./build/emoji.app
sh run.sh build
```

#### Snapshot

![emoji](https://raw.githubusercontent.com/evercyan/cantor/master/resource/47/47deaf5989dd25f6d64bbf4eb6c59f1a.jpg)