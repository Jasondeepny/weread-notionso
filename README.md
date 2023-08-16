# weread-notionso
# golang sync weread note to notion

微信读书笔记、划线等信息同步到notion数据库
> 效果如下：
> ![同步效果](https://markdown-mac-work-1306720256.cos.ap-guangzhou.myqcloud.com/png/AzRZUp.png)

## 使用方法
 ```bash
   1. 获取 **Notion token**
      - 打开[此页面](https://www.notion.so/my-integrations)并登录
      - 点击New integration 输入 name 提交.(如果已有,则点击 view integration)
      - 点击show,然后copy
   2. 从微信读书中获取 cookie
      - 在浏览器中打开 weread.qq.com 并登录
      - 打开开发者工具(按 F12),点击 network(网络),刷新页面, 点击第一个请求,复制 cookie 的值.
   3. 准备 Noiton Database ID
      - 复制[此页面](https://www.notion.so/yayya/d92bb4b8434745baa2061caf67d6ef7a?v=b4a5bfb89e8e44868a473179ee608851)到你的
        Notion中,点击右上角的分享按钮,将页面分享为公开页面
      - 点击页面右上角三个点,在 connections 中找到选择你的 connections.第一步中创建的 integration 的 name
      - 通过 URL 找到你的 Database ID 的值.
      例如:
        页面 https://www.notion.so/yayya/d92bb4b8434745baa2061caf67d6ef7a?v=b4a5bfb89e8e44868a473179ee60x851的ID为d92bb4b8434745baa2061caf67d6ef7a
```

### 1.github Action 运行
 ```bash
    fork 项目
    github action settings >>> secrets >>> 配置环境变量
 ```

### 2.docker 运行
 ```bash
    docker run -d --name weread-notion jasondeepny/weread-notionso:latest
 ```

### 3.docker-compose 运行
 ```bash
    mkdir weread && cd weread
    
    wget https://raw.githubusercontent.com/Jasondeepny/weread-notionso/main/docker-compose.yml
    #修改.env文件环境变量
        wget https://raw.githubusercontent.com/Jasondeepny/weread-notionso/main/.env
    #也可以通过export直接导入环境变量
    
    docker-compose up -d
    
    #加入定定时任务
    crontab -e 
    0 1 * * * cd /root/weread && /usr/local/bin/docker-compose up -d
 ```

## 鸣谢
- https://github.com/malinkang/weread_to_notion

- 配合 NoitonNext 构建 Blog [效果](https://yaya.run/article/1c51da0f-757a-47e5-8296-cc37798b8211)非常好

## 免责申明
本工具仅作技术研究之用，请勿用于商业或违法用途，由于使用该工具导致的侵权或其它问题，该本工具不承担任何责任！

