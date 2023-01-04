var maxHeight = 16056577
var minHeight = 11056577
var amount = parseInt(process.argv[2], 10) || 100
const addr = "http://localhost:8080"

var content = ""
while (amount--) {
    height = getRandomInt(minHeight, maxHeight)
    content += `${addr}/v1/block/0x${height.toString(16)}\n`
}

console.log(content)

// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Math/random
function getRandomInt(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min) + min); // The maximum is exclusive and the minimum is inclusive
}

