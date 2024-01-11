// openssl genrsa -out privateKey.pem 2048
// openssl rsa -in privateKey.pem -pubout -out publicKey.pem

const crypto            = require('crypto')
const fs                = require('fs')
const publicKey         = fs.readFileSync("publicKey.pem")
const privateKey        = fs.readFileSync("privateKey.pem")

const orderNumber               = "123456"

console.log("订单orderNumber>>",orderNumber);

let orderDig         = crypto.createHash('sha256').update(orderNumber,'utf8').digest('hex')
console.log("orderDig:订单orderNumber sha256后结果>>",orderDig);

const signer                       = crypto.createSign('sha256')
signer.update(orderDig)
const sign     = signer.sign(privateKey,'base64')
console.log("签名后的sign>>", sign);
console.log("签名后的sign-encode>>",encodeURIComponent(sign));
const verifier             = crypto.createVerify('sha256')
verifier.update(orderDig);
const ver = verifier.verify(publicKey, sign,'base64');
const validateCode   = encodeURIComponent(sign);
console.log("验签是否成功>>",ver);

fs.writeFileSync("orderDig.txt",orderDig)
fs.writeFileSync("validateCode.txt",validateCode)