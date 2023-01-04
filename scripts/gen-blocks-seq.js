var height = 16056577
var amount = parseInt(process.argv[2], 10) || 100
const addr = "http://localhost:8080"

var content = ""
while (amount--) {
    content += `${addr}/v1/block/0x${height.toString(16)}\n`
    height--
}

console.log(content)