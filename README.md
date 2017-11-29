## 環境変数にほしいもの
- twitter: CONSUMER_KEY, CONSUMER_SECRET, ACCESS_TOKEN, ACCESS_TOKEN_SECRET
- LINE Notify: LINE_AUTHORIZATION

## やりたいこと
- [x] tweetをCSVに保存する
- [x] 指定した文字がstreamに来たらLINE Notifyで通知を送る
- [x] 実行時オプションで通知の動作を変える
- [ ] LINE Notify以外にも通知の方法を選べるようにする
- [ ] ハッシュタグなどEntitiesの要素を入れる
- [ ] streamをgoroutineにする

## ここまでやったけど
ここまでやったけど機能を２つに分割したい
1. UserStreamからTweetの必要部分をCSVに書き出すだけ
2. CSVを監視して標準出力に出す(tail -fを整形)&LINEなどに通知を飛ばす