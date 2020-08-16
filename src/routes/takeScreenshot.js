const puppeteer = require("puppeteer");
const random = require("../util/random");
const ip = require("../util/ip");
const fs = require("fs");

let ipv4 = "";
let ipv6 = "";

module.exports = {
    init: async function() {
        await ip.getIPv4().then(function (result) {
           ipv4 = result;
        });

        await ip.getIPv6().then(function (result) {
            ipv6 = result;
        })

        console.log(`IPv4: ${ipv4}`);
        console.log(`IPv6: ${ipv6}`);

    },

    takeScreenshot: async function(req, repl) {
        try {
            const browser = await puppeteer.launch({
                args: ['--no-sandbox']
            });

            const fileName = random.generateRandomString(8);
            const requestBody = req.body;

            const page = await browser.newPage();

            await page.setViewport({
                width: 1920,
                height: 1080
            });

            await page.goto(requestBody.url);
            await page.waitForResponse(response => response.ok());

            const pageHTML = await page.content();

            if (pageHTML.includes(ipv4)) {
                repl.type("application/json").code(403);
                return {error: "The server's IP address was detected in the response body of the website"};
            }

            if (ipv6 !== undefined) {
                if (pageHTML.includes(ipv6)) {
                    repl.type("application/json").code(403);
                    return {error: "The server's IP address was detected in the response body of the website"};
                }
            }

            await page.screenshot({path: `${fileName}.png`});
            await browser.close()

            const stream = fs.createReadStream(`${fileName}.png`);

            repl.type('image/png');
            repl.send(stream);
            fs.unlinkSync(`${fileName}.png`);
        } catch (err) {
            repl.type('application/json').code(500);
            return {error: err.message};
        }
    }
};