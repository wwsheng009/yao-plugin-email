# 邮件发送插件

yao-plugin-email

email client plugin for yao application。use to send or receive the email。

the default plugin folder path is `<YAO_EXTENSION_ROOT>/plugins/`, the default value for YAO_EXTENSION_ROOT is the app folder, you can change the YAO_EXTENSION_ROOT in the .env file。

## feature

- send email,including the attachement,and send to multi users
- receive email

## build

```sh

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
