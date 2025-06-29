# TLS
## 用語
- 秘密鍵(Private Key)
    - サーバー側で生成して外部に漏らしてはいけない
    - 公開鍵によって暗号化された情報の復号に使う
- 証明書署名要求(CSR)
    - 証明書発行のために作成するファイル
    - サーバの情報を認証局に伝えるための申請書でサーバの正当性を確認するために必要
    - サーバの公開鍵、識別情報（CN, 組織名, 国名など）、署名アルゴリズムなどが含まれる
    - CSRの生成時、サーバはまず秘密鍵と公開鍵のペアを作成
    - 生成した公開鍵とサーバ情報をまとめ、秘密鍵でデジタル署名を行い、CSRを完成させる
- 証明書(Certificate)
    - 認証局(CA)から発行される公開証明書
    - クライアントに対してサーバの正当性を証明できる
- 中間CA証明書(Intermediate Certificate)
    - 信頼チェーンを構成する中間機関の証明書
    - 基本的に中間CAが実際の証明書発行業務を担っている
- ルート証明書(Root Certificate)
    - OSやブラウザに事前にインストールされており，最も信頼されている証明書
    - 認証局（CA）の最上位に位置する証明書
    - クライアント側が管理しているためサーバは関与しない


## TLS設定流れ
1. 秘密鍵の生成
    - サーバー上で秘密鍵（private key）を作成。
    - 例：`openssl genrsa 2048 > server.key`
2. 証明書署名要求（CSR）の作成 ￼
    - 秘密鍵を使い、CSR（Certificate Signing Request）を作成。
    - 例：`openssl req -new -key server.key -out server.csr`
    - CSRにはサーバーの情報（ドメイン名、組織名など）が含まれる。
3. 認証局（CA）へCSRを提出 ￼
    - 作成したCSRを認証局（CA）に送信。
    - CAは申請内容を審査し、問題なければ証明書を発行。
4. サーバー証明書の受領 ￼
    - CAからサーバー証明書（server.crtなど）が発行される。
    - 必要に応じて中間CA証明書も受け取る。
5. サーバーへの証明書設置 ￼
    - サーバーに秘密鍵、サーバー証明書、中間CA証明書を配置。
    - 例：`/etc/ssl/private/server.key`, `/etc/ssl/certs/server.crt`, `/etc/ssl/certs/ca-bundle.crt`
6. サーバー設定ファイルの編集 ￼
    - Webサーバー（Apache, Nginxなど）の設定ファイルで証明書と秘密鍵のパスを指定。
    - 例（Nginx）:￼
7. サーバーの再起動 ￼
    - 設定を反映させるため、Webサーバーを再起動。
8. 動作確認 ￼
    - ブラウザでサイトにアクセスし、証明書が正しく設定されているか確認。
    - SSL Labsなどのツールで検証も可能。





## TLS認証流れ
- `http1-1/README.md`の`### プロセス`を要参照




## 誤解してた点
### その1: サーバに設定される証明書
- サーバには，以下の2点が設定されている
    - 「平文」のサーバ証明書
    - サーバ証明書のハッシュ値をCAの秘密鍵で暗号化したもの

- クライアントはサーバ証明書を取得
- サーバ証明書の「署名部分」を，中間CA証明書の「公開鍵」で検証
- 具体的には、署名部分を中間CAの公開鍵で「復号」し、証明書本体のハッシュ値と一致するかを確認
- これにより、証明書が中間CAによって正しく発行されたこと、改ざんされていないことが証明される


### その2: ルート証明書
- ルート証明書はOSやブラウザに事前にインストールされており，サーバではなくクライアント側で管理している
- 必ず信頼するが，信頼ストアから証明書が削除された場合に信頼しないと判断する




## TLS実装
### 1. 秘密鍵と証明書署名要求(CSR)の作成
```
# 秘密鍵の生成
$ openssl genrsa 2048 > server.key

# 公開鍵の生成(option)
$ openssl rsa -in server.key -pubout -out server-public.key

# 証明書署名要求(csr)を作成
$ openssl req -new -key server.key -subj "/CN=rootca" > server.csr
```


### 2. 証明書の発行
- 本来であればCSRをLet's encryptなどに送付して認証局の秘密鍵で署名してもらう
- 今回は自分で署名することでサーバ証明書を作成するので先ほど生成した秘密鍵で署名する

#### 2-1. サーバの秘密鍵を用いて自己署名
```
# 自己署名証明書の作成
$ openssl x509 -req -in server.csr -days 365 -signkey server.key -out server.crt
```



#### 2-2. 自作のルート認証局を作成して自己署名
**ルート認証局証明書**
```
# 認証局秘密鍵を作成
$ openssl genpkey -algorithm ec -pkeyopt ec_paramgen_curve:prime256v1 -out ca.key

# 証明書署名要求(CSR)を作成
$ openssl req -new -sha256 -key ca.key -out ca.csr -config ca-openssl.cnf

# 証明書を自分の秘密鍵で署名して作成
$ openssl x509 -req -in ca.csr -days 365 -signkey ca.key -sha256 -out ca.crt -extfile ./ca-openssl.cnf -extensions CA
```

**サーバ用**
```
# サーバ秘密鍵を作成
$ openssl genpkey -algorithm ec -pkeyopt ec_paramgen_curve:prime256v1 -out server.key

# 自己署名証明書
$ openssl req -new -nodes -sha256 -key server.key -out server.csr -config server-openssl.cnf

# 証明書を自分の秘密鍵で署名して作成
$ openssl x509 -req -in server.csr -days 365 -sha256 -out server.crt -CA ca.crt -CAkey ca.key -CAcreateserial -extfile ./server-openssl.cnf -extensions Server
```


> Challenge passwordは証明書を破棄するときに使うパスワードだが，自己署名では不要
> 認証局によっては必要になるがセキュリティが弱くあまり使われてない


### 3. HTTPSサーバと証明書の登録
```
# 自己署名CAのルート証明書
$ curl --http1.1 \
       --cacert ca.crt \
       --resolve localhost:8080:127.0.0.1 \
       "https://localhost:8080"
# 明示的にHTTP/1.1を指定
# ローカルの自作CA証明書を使用
# DNSを使わず、名前解決を手動で指定．<ホスト名>:<ポート>:<IPアドレス>
# https を使って localhost で動いている port 8080のTLS対応サーバーにアクセス
# "https://localhost:8080?file=test.txt"と指定するとファイル名を指定できる
```

**Wiresharkで確認**
1. lo0を選択
2. `tcp.port == 8080`を指定


**ルート認証局に自作CAを認識させる**
> [!WARN]
> 完全に自己責任で，攻撃の危険性があるためやらない方がいい！
> 流れだけ確認する

1. openssl x509 -in ca.crt -text -noout
    - 作成した証明書がpem形式か確認
2. cp /opt/homebrew/etc/openssl@3/cert.pem /opt/homebrew/etc/openssl@3/cert.pem.bak
    - バックアップを取る
3. cat ca.crt >> /opt/homebrew/etc/openssl@3/cert.pem
    - 末尾に追記する
4. curl --http1.1 --resolve localhost:8080:127.0.0.1 "https://localhost:8080"
    - ルートCA証明書を指定しなくても通信できることを確認
    - すでに信頼されているルート証明書のリストに書き込んだから
5. cp /opt/homebrew/etc/openssl@3/cert.pem.bak /opt/homebrew/etc/openssl@3/cert.pem
    - 信頼されているルート証明書をもとに戻す


### Client
```
$ openssl genpkey -algorithm ec -pkeyopt ec_paramgen_curve:prime256v1 -out client.key

$ openssl req -new -nodes -sha256 -key client.key -out client.csr -config client-openssl.cnf

$ openssl x509 -req -in client.csr -days 365 -sha256 -out client.crt -CA ca.crt -CAkey ca.key -CAcreateserial -extfile ./client-openssl.cnf -extensions Client
```

- サーバがクライアントの証明書を要求した際に送信する
- 証明書と同時に公開鍵も送られるため，それらをもとにサーバは検証を行う

#### mTLS(相互TLS認証)
- サーバ，クライアントの両方が証明書を持って信頼する相手のみと通信すること
- ウェブブラウジングでは利用されず，高いセキュリティが求められるときに使われる
    - VPN，IoT，APIなど
- 信頼されたクライアントからの通信のみ受け付ける場合

