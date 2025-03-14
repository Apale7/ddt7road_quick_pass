# 弹弹堂七道经典服糖果面验证码登录插件
## 使用
1. 下载并解压插件到本地，双击启动
2. 糖果浏览器侧边栏的账号设置，每个账号需要右键-编辑账号，在游戏网址后面加上账号密码两个参数，格式为
username=xxx&pwd=yyy
其中xxx为账号，yyy为密码明文或密码的base64编码(糖果导出的是base64)
1. 糖果右上角 工具-代理服务器设置-管理代理服务器列表 添加，地址：127.0.0.1，端口号：8888，名称：七道经典服免验证登录，类型：http。添加后确认已启用该代理
## 原理
弹弹堂官网点击登录按钮时的流程存在漏洞
1. 跳转到滑块验证码
2. 验证码通过后跳转到账号密码验证
3. 验证通过再跳转到游戏页面
问题出现在第二步：账号密码验证的接口并没有关注是否通过了滑块验证码
因此直接在后台登录并写入cookie即可跳过验证码完成登录

## 插件登录流程
1. 拦截http://www.wan.com/game/play/id/8665.html?username=aaa&pwd=bbb的response
2. 从请求的query中获取账号密码
3. 后台调用登录接口，获得cookies
4. cookies写入response的Set-Cookie header
此时因为response已经包含了登录后的cookies，所以已经登录成功，验证码不会弹出，直接进入游戏加载界面

## Q&A
1. 插件登录时怎么获取到的账号密码？
- A: 插件无法获取糖果保存的账号密码，因此需要手动把账号密码写到糖果游戏网址的query参数中。如果号太多建议先导出糖果账号，然后批量修改账号文件后再导入，会方便一点![](http://qiniu.apale7.cn/20250308200927.png)

  图中蓝色部分是需要新增的
2. 官方会检测到吗？会封号吗？
- A: 官方可以检测到这次登录没有经过验证码。不确定会不会封号，但目前另一款常用登录器36脚本大厅免验证登录也是类似的原理，应该不会被封