const request = require("then-request");

module.exports = {
    getIPv4: function () {
        return new Promise(function (resolve, reject) {
            request("GET", "http://ip-api.com/json/")
                .getBody("utf-8")
                .then(JSON.parse)
                .done(function (res) {
                    resolve(res.query);
                });
        });
    },

    getIPv6: function () {
        return new Promise(function (resolve, reject) {
            request("GET", "https://api6.ipify.org?format=json", {family: 6})
                .getBody("utf-8")
                .then(JSON.parse)
                .done(function (res) {
                    resolve(res.ip);
                });
        })
    }
}