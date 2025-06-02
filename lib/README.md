# ReeX Library

## 機能

agentの起動

- push型
- 設定ファイルをもとにagentを起動
  - nodename IP User "ssh-pass:xxxx or pubkey:/path/to/pubkey" "/path/to/your/AgentDir"
  - sudo を使用したければ最初からrootユーザを指定
- 初回はssh接続を用いたagentのコピー・起動
  - "ps aux"を確認し，既に起動しているかを確認
  - agentをコピー
  - バックグラウンドでagentを起動(sudo)
  - agentのkillコマンド(api)を用意

---

コマンドの実行

- コピーはsshを利用
- 設定ファイルを用いてセッションを作成
  - nodename IP Port "/path/to/your/workdir"  
    //

## 構成

```
.
└── lib/
    ├── controller/
    │   ├── ssh
    │   ├── exec
    │   ├── api/
    │   │   ├── ssh
    │   │   └── session
    │   └── config/
    │       ├── confssh
    │       └── confsession
    └── agent
```
