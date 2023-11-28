# 邮件发送插件

yao-plugin-email

[Document](https://wwsheng009.github.io/yao-docs/YaoDSL/Plugin/golang%20grpc%20%E6%8F%92%E4%BB%B6%E6%A8%A1%E6%9D%BF.html)

email send plugin for yao application。

the default plugin folder path is `<YAO_EXTENSION_ROOT>/plugins/`, the default value for YAO_EXTENSION_ROOT is the app folder, you can change the YAO_EXTENSION_ROOT in the .env file。

## build

```sh

cd plugins/email

make build
```

## test

maintain the .env file

```sh
EMAIL_USERNAME=
EMAIL_PASSWORD=
EMAIL_TO=
```

run test use yao command

```js
yao run scripts.email.test
```
