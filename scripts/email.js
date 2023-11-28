
//yao run scripts.email.test
function test() {
    const username = Process("utils.env.Get", "EMAIL_USERNAME")
    const to = Process("utils.env.Get", "EMAIL_TO")
    const password = Process("utils.env.Get", "EMAIL_PASSWORD")

    const message = {
        "account": {
            "server": "smtp.qq.com",
            "port": 587,
            "username": username,
            "password": password,
            "type": "stmp"
        },
        "messages": [
            {
                "from": username,
                "to": to,
                "cc": [username],
                "subject": "小佩奇",
                "body": "<h1>新年快乐</h1>",
                "attachments": ["./data/test.jpg"]
            }
        ]
    }
    const resp = Process("plugins.email.send", message)
    console.log(resp)
}