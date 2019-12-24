
const fs = require('fs')
const dataPath = './data.json'

async function readFile() {
  return new Promise((resolve, reject) => {
    fs.readFile(dataPath, 'utf8', (err, result) => {
      if (err) {
        console.log(err)
        reject(err)
      }

      const data = JSON.parse(result)
      resolve(data)
    })
  })
}

const NodeRSA = require('node-rsa');

const _pubKey = `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAICZifH6EG/qZ4tNYdkTUI/+TMpMr/3r
ap6NM5Zqo0pEaTQWRVarR3njktB3ssfydQZESo+E0d3AIffhXrf1YiUCAwEAAQ==
-----END PUBLIC KEY-----`;

const _priKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAICZifH6EG/qZ4tNYdkTUI/+TMpMr/3rap6NM5Zqo0pEaTQWRVar
R3njktB3ssfydQZESo+E0d3AIffhXrf1YiUCAwEAAQJAce6VZgoTsfNMFQBTpgwU
nd3AjqHucdm1tg6HG8YSMOKnsKi5YWr1nB7VOGgbN6UVIqMCphfbkglIA2A2ICQP
AQIhAPMVXuPOysQdI0qQ0gD5t79shLr5TSCnsQDN4QuN6CaFAiEAh27c3ebGy3s1
AtggAQvyQYUB4Gi/tr2ZamwlMh7+LyECIQCW/IgzEehKRhr8ntWCO5nRcdND28P3
a5F7AWYuahdvjQIgewYNw9S6iGRnBypkCA9eBH5Z8gu0+r7H+ZA7SYg1xYECIDKd
p6W2H0D+k/G4yazk0SSfNg7fNCZcdY8qMwTfciOO
-----END RSA PRIVATE KEY-----`;

function encryption(data) {

  const publicKey = new NodeRSA(_pubKey)
  encrypted = publicKey.encrypt(data)
  // console.log('encrypted data code:', encrypted)

  return encrypted
}

function decryption(code) {

  const privateKey = new NodeRSA(_priKey)
  decrypted = privateKey.decrypt(code, 'utf8')
  // console.log('decrypted data:', decrypted)

  return decrypted
}

const path = require('path')
const mimeType = require('mime-types')

function parseBase64(img) {

  const imgPath = path.resolve(img)
  const imgMimeType = mimeType.lookup(imgPath)
  const data = fs.readFileSync(imgPath)
  const dataBase64 = `data:${imgMimeType};base64,${Buffer.from(data, 'base64').toString('base64')}`
  // console.log('data base64:', dataBase64)

  return dataBase64

}

async function parseFile(name, dataBase64) {
  const dataBuffStr = dataBase64.replace(/^data:image\/\w+;base64,/, '')
  const dataBuff = Buffer.from(dataBuffStr, 'base64')

  return new Promise((resolve, reject) => {
    fs.writeFile(name + '.jpg', dataBuff, err => {
      if (err) {
        console.log('save file err:', err)
        reject(err)
      } else {
        resolve()
      }
    })
  })
}



const { query } = require('./query')
const { invoke } = require('./invoke')

const MAIN_CHANNEL = {
  channel: 'mainchannel',
  chaincode: 'mainchannel'
}

const SUB_CHANNEL = {
  channel: 'subchannel',
  chaincode: 'subchannel'
}

function translationToArr(obj) {
  const keys = ['user', 'channel', 'chaincode', 'fcn', 'args']
  return keys.map(val => {
    if (obj[val]) {
      return obj[val]
    }
  })
}

function dateFormat(date) {
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()

  const hour = date.getHours()
  const minute = date.getMinutes()
  const second = date.getSeconds()

  return {
    full: `${year}-${month}-${day} ${hour}:${minute}:${second}`,
    time: date.getTime()
  }
}

function sleep(time) {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve()
    }, time)
  })
}

async function main() {

  const costArr = []

  const res = await readFile()
  const data = res.slice(0, 3500)
  console.log('data length:', data.length)

  const len = data.length

  for(let i=0; i<len; i++) {
    const dataVal = data[i]
    const { channel, chaincode } = SUB_CHANNEL

    dataVal['photo'] = encryption(parseBase64(dataVal['photo_path'])).toString('hex')

    // 写入
    const date = new Date()
    const dateObj = dateFormat(date)
    console.log('====  start  ====')

    try {
      const result = await invoke(
        'user1',
        channel,
        chaincode,
        'dataUpload',
        [
          'testRequest',
          dataVal['photo_id'],
          dataVal['photo'],
          dateObj['full']
        ]
      )

      if (!result.success) {
        continue
      }

      const cost = (new Date().getTime() - dateObj.time) / 1000
      costArr.push({
        name: dataVal['photo_id'],
        path: dataVal['photo_path'],
        cost
      })

      console.log('cost:', cost)
      console.log('count:', costArr.length)
    } catch(e) {
      console.log('err:', e)

      return
    }

    console.log('====  end  =====')


    // // 读取
    // const result = await query(
    //   PEER0_ORG0_TRUST_COM,
    //   SUB_CHANNEL,
    //   'dataQuery',
    //   [
    //     'testRequest',
    //     dataVal['photo_id']
    //   ]
    // )
    //
    //
    // // 生成文件
    // await parseFile(result.name, decryption(Buffer.from(result.content, 'hex')))
  }

  await fs.writeFile('./cost.json', JSON.stringify(costArr), err => {
    if (err) {
      console.log('write file err:', err)
    }
  })
}

main()
