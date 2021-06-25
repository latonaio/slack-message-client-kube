# slack-message-client-kube
kanbanで受け取ったデータをslackに通知するマイクロサービスです。

## 動作環境
動作には以下の環境であることを前提とします。

- OS: Linux
  
- CPU: Intel64/AMD64/ARM64

最低限スペック  
- CPU: 2 core  
  
- memory: 4 GB

## 起動方法
docker imageのビルド
```
$ git clone {slack-message-client-kube}
$ cd ~/path/to/slack-message-client-kube
$ bash docker-build.sh
```

project.yamlに次の設定を追加してください。
```yaml
slack-message-client-kube:
startup: yes
always: yes
env:
  TOKEN: XXXXX
  CHANNEL_ID: XXX
```
  

## 環境変数
- CHANNEL_ID: 通知先チャンネルのID
- TOKEN: slack apiのOAuthトークン

## Input  
kanbanデータを受け取ります。
必要パラメータ：
```
- pod_name string
- status string
```
  
## Output  
受け取ったデータを整形してslackに通知します。
メッセージの大量送信を防ぐために、同一podからのメッセージは5回までとなっています。
メッセージのlevelが"warning"のものに限り、slackにメッセージを送信します。

## slack連携方法
このマイクロサービスの実行にはslack WEB apiの利用設定と発行されるOAuth TOKENが必要です。
詳細はこちらを確認してください。
- https://api.slack.com/web