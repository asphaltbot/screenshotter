const puppeteer = require("puppeteer");
const random = require("../util/random");
const ip = require("../util/ip");
const fs = require("fs");

let browser = undefined;

module.exports = {
    init: async function() {
        browser = await puppeteer.launch({
            args: ['--no-sandbox']
        });
    },

    takeScreenshot: async function(req, repl) {
        try {
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

            let wasIPDetected = false;
            await ip.getIP().then(function (result) {
                if (pageHTML.includes(result)) {
                    wasIPDetected = true;
                }
            });

            if (wasIPDetected) {
                repl.type("application/json").code(403);
                return {error: "The server's IP address was detected in the response body of the website"};
            }

            await page.screenshot({path: `${fileName}.png`});

            const stream = fs.createReadStream(`${fileName}.png`);
            fs.unlinkSync(`${fileName}.png`);

            repl.type('image/png');
            repl.send(stream);
        } catch (err) {
            repl.type('application/json').code(500);
            return {error: err.message};
        }
    }
};