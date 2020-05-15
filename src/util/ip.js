const request = require("then-request");

module.exports = {
    getIP: function () {
        return new Promise(function (resolve, reject) {
            request("GET", "http://ip-api.com/json/")
                .getBody("utf-8")
                .then(JSON.parse)
                .done(function (res) {
                    resolve(res.query);
                });
        });
    }
}