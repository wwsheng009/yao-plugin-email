
//yao run scripts.email.send
function send() {
    const username = Process("utils.env.Get", "EMAIL_USERNAME")
    const to = Process("utils.env.Get", "EMAIL_TO")
    const password = Process("utils.env.Get", "EMAIL_PASSWORD")

    // you can send to multi users
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
                "to": [
                    { name: 'vincent', address: to },
                    { name: 'vincent1', address: to }
                  ],
                "cc": [{ name: 'myown', address: to }],
                "subject": "小佩奇",
                "body": "<h1>新年快乐</h1>",
                "attachments": ["./data/test.jpg"]
            }
        ]
    }
    const resp = Process("plugins.email.send", message)
    console.log(resp)
}

//yao run scripts.email.receive
function receive() {
    const username = Process("utils.env.Get", "EMAIL_USERNAME")
    const to = Process("utils.env.Get", "EMAIL_TO")
    const password = Process("utils.env.Get", "EMAIL_PASSWORD")

    const message = {
        "account": {
            "server": "imap.qq.com",
            "port": 993,
            "username": username,
            "password": password,
            "type": "imap"
        }
    }
    const resp = Process("plugins.email.receive", message)
    console.log(resp)
}