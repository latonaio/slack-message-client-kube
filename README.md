# slack-message-client-kube
エッジのメッセージデータをslackに通知するマイクロサービスです。


## 起動方法
docker imageのビルド
```
$ cd ~/path/to/slack-message-client-kube
$ bash docker-build.sh
```

`slack-message-client-kube.yml` に次の設定を追加してください。

`QUEUE_FROM` で指定するキューは事前に作成が必要です。

```yaml
env:
- name: CHANNEL_ID
  value: "チャンネル ID"
- name: TOKEN
  value: "Slack App トークン"
- name: RABBITMQ_URL
  value: "amqp://username:password@hostname:5672/virtualhost"
- name: QUEUE_FROM
  value: "キュー名"
```
  
## 動作環境
動作には以下の環境であることを前提とします。

```
- OS: Linux
- CPU: Intel64/AMD64/ARM64

最低限スペック  
- CPU: 2 core  
- memory: 4 GB
```

## 環境変数
- CHANNEL_ID: 通知先チャンネルのID
- TOKEN: slack apiのOAuthトークン
- RABBITMQ_URL: RabbitMQのURL
- QUEUE_FROM: RabbitMQの受信元キュー名

## Input  
メッセージデータを受け取ります。
必要パラメータ：
```
- pod_name string
- status string
- level string
```
  
## Output  
受け取ったデータを整形してslackに通知します。
メッセージのlevelが"warning"のものに限り、slackにメッセージを送信します。

## slack連携方法
このマイクロサービスの実行にはslack WEB apiの利用設定と発行されるOAuth TOKENが必要です。
詳細はこちらを確認してください。
- https://api.slack.com/web