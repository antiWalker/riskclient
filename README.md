# 风控规则引擎

[风控规则引擎 Go 语言版接口文档](https://wiki.mafengwo.cn/pages/viewpage.action?pageId=29925339)

### 部署到线上
修改 `Dockerfile` 和 `pipeline.yaml`. *按照文件中的说明修改*

### 接口模拟测试
RiskEngine.json 是 [Insomnia](https://insomnia.rest/) 导出的 json. 
下载  Insomnia 再导入之后可以直接发送模拟请求.

### 自动生成model的使用步骤
1 本地数据库需要有你要依赖生成model的表结构。<br>
2 执行命令：./tools/models-gen-go -dsn "root:123456@tcp(localhost:3306)/im?charset=utf8" -file "/Users/wang/salesorder.go" -table "t_risk_engine_sales_order" -structName "SalesOrder" <br>
3 把依赖的表所生成的xxx.go文件拷贝到models文件夹下 <br>
