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

    getIPv6: function() {
        return new Promise(function (resolve, reject) {
            try {
                request("GET", "http://v6.ipv6-test.com/api/myip.php")
                    .getBody("utf-8")
                    .error(e => console.log(e))
                    .done(function (res) {
                        resolve(res);
                    })

            } catch (e) {
                resolve(undefined);
            }
        })
    }

}