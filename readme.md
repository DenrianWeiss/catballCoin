# CatballCoin

"前后端分离式虚拟币诈骗"  
本程序的作用是通过自动查询 RPC 获取账户余额并尽可能快的转账出eth

## 使用

1. 修改配置文件，将 `config.example.json` 移动到 `config.json`，并修改 RPC 密钥。
2. 调用 rpc 添加任务

## RPC

向以下地址发送 post 请求：`/task/add/您的API密钥`。请求体

```json
{
  "target_wallet": "你的钱包地址",
  "source_key": "源钱包私钥",
  "chat_id": 机器人要发送消息的telegram chat id,
  "gas_mode": "费用计算方式，看下文",
  "eth_rpc": "以太坊 RPC 端点",
  "low_threshold": 最低触发交易的余额 wei 数,
  "high_threshold": gas总费用上线，wei,
  "percents": 计算gas的参数(百分比乘以100)
}
```

其中 `gas_mode` 可以是 `percent` 或 `price`，前者则总 gas 费用为金额乘百分比，后者则 gas 费用为当前网络 gas 费用的 percent 倍
