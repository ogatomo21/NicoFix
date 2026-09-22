# NicoFix

Discordのメッセージに含まれるニコニコ動画URLを検出し、`nicovideo.jp` を `nicovideo.gay` に置き換えたURLをReplyで返すBotです。

## Botをインストール

[DiscordにBotをインストールする](https://discord.com/oauth2/authorize?client_id=1551980109341130903)

対象は `http(s)://nicovideo.jp/watch/sm数字` と `http(s)://www.nicovideo.jp/watch/sm数字` です。クエリ文字列とフラグメントは維持し、同じURLは1回だけ返信します。

## Credit

本Botは、MMakerさんの [nicovideo.gay](https://mmaker.moe/2025/01/niconico-geo-block-user-agent-spoofing/) へのURL自動置換を行います。

## Discord Botの作成

1. [Discord Developer Portal](https://discord.com/developers/applications) で **New Application** を選び、任意の名前で作成します。
2. 左メニューの **Bot** を開き、**Reset Token** でBot Tokenを発行して控えます。Tokenは公開しないでください。
3. 同じ **Bot** ページの **Privileged Gateway Intents** で **Message Content Intent** を有効にして保存します。
4. **OAuth2 > URL Generator** で Scope に `bot` を選びます。Bot Permissions は次の最小権限を選択します。
   - `View Channels`
   - `Send Messages`
   - `Manage Messages`
5. 生成されたURLを開き、Botを対象サーバーへ招待します。

Message Content Intentが無効だと、Botはメッセージ本文を取得できずURLを変換できません。

## 自己ホスト

公開BotはTomoya Ogawa(ogatomo21)が運営していますが、継続提供は保証しません。安定して使いたい場合は、ローカルPCまたは任意のLinuxサーバーで自身のBotをホストしてください。公開Bot・自己ホスト版ともに利用は自己責任です。

## ローカル実行

PowerShellでプロジェクト直下を開き、Tokenを現在のセッションだけに設定して実行します。

```powershell
$env:DISCORD_TOKEN = "Bot Token"
go run .
```

`DISCORD_TOKEN` が未設定の場合は、明確なエラーを表示して終了します。

## ビルド

```powershell
go build -o NicoFix.exe .
```

Linux amd64向けのリリースビルド:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o NicoFix .
```

テスト・静的検査:

```bash
go test ./...
go vet ./...
```

## GCE / Debianへの配置

ローカルPCでLinux向けにビルドした `NicoFix` と `nicofix.service` をGCEへコピーしてから、GCE上で実行します。

```bash
sudo install -d -m 0755 /opt/NicoFix
sudo install -m 0755 NicoFix /opt/NicoFix/NicoFix
sudo install -m 0644 nicofix.service /etc/systemd/system/nicofix.service
sudo useradd --system --user-group --home-dir /opt/NicoFix --shell /usr/sbin/nologin nicofix
sudo install -m 0600 /dev/null /etc/nicofix.env
sudoedit /etc/nicofix.env
```

`/etc/nicofix.env` の内容:

```ini
DISCORD_TOKEN=Bot Token
```

Tokenを保存したら、サービスを有効化して起動します。

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now nicofix
```

## systemd操作

```bash
sudo systemctl status nicofix
sudo systemctl restart nicofix
sudo systemctl stop nicofix
sudo journalctl -u nicofix -f
```

## ライセンス

MIT License
