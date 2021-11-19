# slack-message-client-kube
エッジのメッセージデータをslackに通知するマイクロサービスです。


## 起動方法
docker imageのビルド
```
$ cd ~/path/to/slack-message-client-kube
$ bash docker-build.sh
```

`slack-message-client-kube.yml` に次の設定を追加してください。

`QUEUE_ORIGIN` で指定するキューは事前に作成が必要です。

```yaml
env:
- name: CHANNEL_ID
  value: "チャンネル ID"
- name: TOKEN
  value: "Slack App トークン"
- name: RABBITMQ_URL
  value: "amqp://username:password@hostname:5672/virtualhost"
- name: QUEUE_ORIGIN
  value: "キュー名"
```
  
## 動作環境
動作には以下の環境であることを前提とします。

```
- OS: Linux
- CPU: ARM/AMD/Intel

最低限スペック  
- CPU: 2 core  
- memory: 4 GB
```

## 環境変数
- CHANNEL_ID: 通知先チャンネルのID
- TOKEN: slack apiのOAuthトークン
- RABBITMQ_URL: RabbitMQのURL
- QUEUE_ORIGIN: RabbitMQの受信元キュー名

## Input  
RabbitMQ からメッセージデータを受け取ります。
必要パラメータ：
```
- pod_name string
- status string
- level string
```
  
入力データのサンプル は、inputs/sample.json にある通り、次の様式です。
```
{
   "terminalName": "xxxxx",
   "macAddress": "xx:xx:xx:xx:xx:xx",
   "createdAt": "2021-10-16T03:13:27.539Z",
   "imagePath": "/var/lib/aion/Data/direct-next-service_1/1634354006794.jpg",
   "faceId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
   "responseData": {
       "candidates": []
   }
}
```
  
## Output  
データを整形してslackに通知します。
メッセージのlevelが"warning"のものに限り、slackにメッセージを送信します。

## slack連携方法
このマイクロサービスの実行にはslack WEB apiの利用設定と発行されるOAuth TOKENが必要です。
詳細はこちらを確認してください。
- https://api.slack.com/web