const https = require('https');  

const defaultOptions = {
    host: 'example.com',
    port: 443,
    headers: {
        'Content-Type': 'application/json',
    }
}

const post = (path, payload) => new Promise((resolve, reject) => {
    const options = { ...defaultOptions, path, method: 'POST' };
    const req = http.request(options, res => {
        let buffer = "";
        res.on('data', chunk => buffer += chunk)
        res.on('end', () => resolve(JSON.parse(buffer)))
    });
    req.on('error', e => reject(e.message));
    req.write(JSON.stringify(payload));
    req.end();
})

// Example usage
exports.handler = async (event, context) => new Promise( async (resolve, reject) => {
    
    const result = await post("/s3/scanfile", { 	bucketname = "hawk-s3-virus-scan-test",
	key = "eicar" });

    console.log(JSON.stringify(result));

})